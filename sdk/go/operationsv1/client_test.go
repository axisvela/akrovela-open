package operationsv1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientSubmit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/operations" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		var request SubmitOperationRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode: %v", err)
		}
		if request.OperationKey != "op-1" {
			t.Errorf("operation = %q", request.OperationKey)
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(OperationView{OperationKey: "op-1", ExecutionStatus: "pending"})
	}))
	defer server.Close()

	view, err := NewClient(server.URL, "token").Submit(context.Background(), SubmitOperationRequest{OperationKey: "op-1"})
	if err != nil {
		t.Fatal(err)
	}
	if view.OperationKey != "op-1" || view.ExecutionStatus != "pending" {
		t.Fatalf("view = %+v", view)
	}
}

func TestClientQueryCancelReceipts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/operations/op-1":
			if r.Method != http.MethodGet {
				t.Errorf("query method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(OperationView{OperationKey: "op-1"})
		case "/v1/operations/op-1/cancel":
			if r.Method != http.MethodPost {
				t.Errorf("cancel method = %s", r.Method)
			}
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(CancelResponse{CancelRequestID: "cancel-1", CancelStatus: "unknown"})
		case "/v1/operations/op-1/receipts":
			if r.URL.Query().Get("after_sequence") != "4" {
				t.Errorf("after_sequence = %q", r.URL.Query().Get("after_sequence"))
			}
			_ = json.NewEncoder(w).Encode(ReceiptList{NextAfterSequence: 5})
		case "/v1/receipts/receipt-1":
			_ = json.NewEncoder(w).Encode(Receipt{ReceiptID: "receipt-1"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	ctx := context.Background()

	view, err := client.Get(ctx, "op-1")
	if err != nil || view.OperationKey != "op-1" {
		t.Fatalf("get = %+v err=%v", view, err)
	}
	cancel, err := client.Cancel(ctx, "op-1", CancelRequest{CancelRequestID: "cancel-1"})
	if err != nil || cancel.CancelRequestID != "cancel-1" {
		t.Fatalf("cancel = %+v err=%v", cancel, err)
	}
	receipts, err := client.Receipts(ctx, "op-1", 4)
	if err != nil || receipts.NextAfterSequence != 5 {
		t.Fatalf("receipts = %+v err=%v", receipts, err)
	}
	receipt, err := client.GetReceipt(ctx, "receipt-1")
	if err != nil || receipt.ReceiptID != "receipt-1" {
		t.Fatalf("receipt = %+v err=%v", receipt, err)
	}
}

func TestClientMapsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(Error{Error: ErrorDetail{
			Code: "operation_key_binding_conflict", Message: "binding differs", Retryable: false,
		}})
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "token").Get(context.Background(), "op-1")
	var apiError *APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("err = %v", err)
	}
	if apiError.StatusCode != http.StatusConflict || apiError.Code != "operation_key_binding_conflict" {
		t.Fatalf("api error = %+v", apiError)
	}
}

func TestClientRequiresBaseURL(t *testing.T) {
	if _, err := NewClient("", "").Get(context.Background(), "op-1"); err == nil {
		t.Fatal("empty base URL was accepted")
	}
}
