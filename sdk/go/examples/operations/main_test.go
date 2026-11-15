package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	operationsv1 "github.com/axisvela/akrovela-open/sdk/go/operationsv1"
)

func TestRunSubmitFillsDigests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/operations" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		var request operationsv1.SubmitOperationRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode: %v", err)
		}
		if !strings.HasPrefix(request.PayloadDigest, "sha256:") ||
			!strings.HasPrefix(request.BindingDigest, "sha256:") {
			t.Errorf("digests = %q %q", request.PayloadDigest, request.BindingDigest)
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(operationsv1.OperationView{OperationKey: "op-1"})
	}))
	defer server.Close()

	path := filepath.Join(t.TempDir(), "request.json")
	request, err := json.Marshal(operationsv1.SubmitOperationRequest{
		OperationKey: "op-1", OperationType: "ros2.nav2.navigate_to_pose@1",
		Scope: operationsv1.Scope{TenantID: "tenant"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, request, 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err = run([]string{"submit", "-request", path, "-base-url", server.URL, "-token", "token"}, strings.NewReader(""), &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"operation_key": "op-1"`) {
		t.Fatalf("output = %s", out.String())
	}
}

func TestRunReadsStdinAndRejectsUnknownCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(operationsv1.OperationView{OperationKey: "op-1"})
	}))
	defer server.Close()
	request, err := json.Marshal(operationsv1.SubmitOperationRequest{OperationKey: "op-1"})
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"submit", "-request", "-", "-base-url", server.URL}, bytes.NewReader(request), &out); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"unknown"}, strings.NewReader(""), &out); err == nil {
		t.Fatal("unknown command was accepted")
	}
}

func TestRunSubmitWait(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/operations":
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(operationsv1.OperationView{OperationKey: "op-1", ExecutionStatus: "started"})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/operations/op-1":
			_ = json.NewEncoder(w).Encode(operationsv1.OperationView{OperationKey: "op-1", ExecutionStatus: "completed"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	path := filepath.Join(t.TempDir(), "request.json")
	request, err := json.Marshal(operationsv1.SubmitOperationRequest{
		OperationKey: "op-1", OperationType: "ros2.nav2.navigate_to_pose@1",
		Scope: operationsv1.Scope{TenantID: "tenant"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, request, 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err = run([]string{"submit", "-request", path, "-base-url", server.URL, "-wait", "-interval-ms", "1"}, strings.NewReader(""), &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"execution_status": "completed"`) {
		t.Fatalf("output = %s", out.String())
	}
}
