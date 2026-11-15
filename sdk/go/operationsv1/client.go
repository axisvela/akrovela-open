package operationsv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Retryable  bool
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("operationsv1: %s (%s)", e.Message, e.Code)
	}
	return fmt.Sprintf("operationsv1: http %d", e.StatusCode)
}

type Client struct {
	BaseURL string
	HTTP    *http.Client
	Token   string
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		HTTP:    &http.Client{Timeout: 10 * time.Second},
		Token:   token,
	}
}

// TerminalExecutionStatus reports whether execution will not advance further.
// Effect status remains separate; unknown execution is not terminal.
func TerminalExecutionStatus(status string) bool {
	switch status {
	case "completed", "failed", "confirmed_no_effect":
		return true
	default:
		return false
	}
}

// Wait polls until execution reaches a terminal status or the context ends. It
// does not wait for effect verification and it does not retry a submitted effect.
func (c *Client) Wait(ctx context.Context, operationKey string, interval time.Duration) (OperationView, error) {
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		view, err := c.Get(ctx, operationKey)
		if err != nil {
			return view, err
		}
		if TerminalExecutionStatus(view.ExecutionStatus) {
			return view, nil
		}
		select {
		case <-ctx.Done():
			return view, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c *Client) Submit(ctx context.Context, request SubmitOperationRequest) (OperationView, error) {
	var view OperationView
	err := c.do(ctx, http.MethodPost, "/v1/operations", nil, request, &view, http.StatusAccepted, http.StatusOK)
	return view, err
}

func (c *Client) Get(ctx context.Context, operationKey string) (OperationView, error) {
	var view OperationView
	err := c.do(ctx, http.MethodGet, "/v1/operations/"+url.PathEscape(operationKey), nil, nil, &view, http.StatusOK)
	return view, err
}

func (c *Client) Cancel(ctx context.Context, operationKey string, request CancelRequest) (CancelResponse, error) {
	var response CancelResponse
	err := c.do(ctx, http.MethodPost, "/v1/operations/"+url.PathEscape(operationKey)+"/cancel", nil, request, &response, http.StatusAccepted)
	return response, err
}

func (c *Client) Receipts(ctx context.Context, operationKey string, afterSequence int64) (ReceiptList, error) {
	query := url.Values{}
	if afterSequence > 0 {
		query.Set("after_sequence", fmt.Sprintf("%d", afterSequence))
	}
	var list ReceiptList
	err := c.do(ctx, http.MethodGet, "/v1/operations/"+url.PathEscape(operationKey)+"/receipts", query, nil, &list, http.StatusOK)
	return list, err
}

func (c *Client) GetReceipt(ctx context.Context, receiptID string) (Receipt, error) {
	var receipt Receipt
	err := c.do(ctx, http.MethodGet, "/v1/receipts/"+url.PathEscape(receiptID), nil, nil, &receipt, http.StatusOK)
	return receipt, err
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any, statuses ...int) error {
	if c.BaseURL == "" {
		return fmt.Errorf("operationsv1: base URL is required")
	}
	address := c.BaseURL + path
	if len(query) > 0 {
		address += "?" + query.Encode()
	}
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(ctx, method, address, reader)
	if err != nil {
		return err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}
	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, 4<<20))
	if err != nil {
		return err
	}
	for _, status := range statuses {
		if response.StatusCode == status {
			if out == nil {
				return nil
			}
			if err := json.Unmarshal(payload, out); err != nil {
				return err
			}
			return nil
		}
	}
	var apiError Error
	if err := json.Unmarshal(payload, &apiError); err == nil && apiError.Error.Code != "" {
		return &APIError{
			StatusCode: response.StatusCode,
			Code:       apiError.Error.Code,
			Message:    apiError.Error.Message,
			Retryable:  apiError.Error.Retryable,
		}
	}
	return &APIError{StatusCode: response.StatusCode}
}
