package akrovela

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	base string
	http *http.Client
}

type VersionInfo struct {
	Product string `json:"product"`
	Version string `json:"version"`
	Stage   string `json:"stage"`
}

func New(base string) (*Client, error) {
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("HTTP(S) base URL without credentials, query or fragment required")
	}
	return &Client{base: strings.TrimRight(base, "/"), http: &http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}, nil
}

func (c *Client) Version(ctx context.Context) (VersionInfo, error) {
	var result VersionInfo
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/version", nil)
	if err != nil {
		return result, errors.New("invalid version request")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}
		return result, errors.New("version request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("version returned HTTP %d", resp.StatusCode)
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 65536)).Decode(&result); err != nil {
		return result, errors.New("invalid version response")
	}
	if result.Product != "AKROVELA" || result.Version == "" || result.Stage != "bootstrap" {
		return VersionInfo{}, errors.New("unexpected product version response")
	}
	return result, nil
}
