package akrovela

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRejectsInvalidURL(t *testing.T) {
	for _, base := range []string{"", "ftp://host", "http://u:p@host", "http://host?q=1", "http://host#fragment"} {
		if _, err := New(base); err == nil {
			t.Fatalf("accepted %q", base)
		}
	}
}

func TestVersion(t *testing.T) {
	for _, tc := range []struct {
		status int
		body   string
		ok     bool
	}{
		{200, `{"product":"AKROVELA","version":"0.1.0","stage":"bootstrap"}`, true},
		{503, `unavailable`, false},
		{200, `{`, false},
		{200, `{"product":"OTHER","version":"0.1.0","stage":"bootstrap"}`, false},
	} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || r.URL.Path != "/version" {
				t.Error("unexpected request")
			}
			w.WriteHeader(tc.status)
			_, _ = w.Write([]byte(tc.body))
		}))
		client, err := New(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		result, err := client.Version(context.Background())
		server.Close()
		if (err == nil) != tc.ok || (tc.ok && result.Version != "0.1.0") {
			t.Fatalf("result=%+v error=%v", result, err)
		}
	}
}

func TestVersionCancellation(t *testing.T) {
	client, _ := New("http://127.0.0.1:1")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Version(ctx); err != context.Canceled {
		t.Fatalf("error=%v", err)
	}
}
