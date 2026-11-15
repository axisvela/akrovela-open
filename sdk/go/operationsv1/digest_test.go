package operationsv1

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

func TestCanonicalNumberVectors(t *testing.T) {
	cases := []struct {
		value    any
		expected string
	}{
		{1e-6, "0.000001"},
		{1.234e-6, "0.000001234"},
		{1e-7, "1e-7"},
		{1e20, "100000000000000000000"},
		{1e21, "1e+21"},
		{0.45, "0.45"},
		{40.0, "40"},
		{math.Copysign(0, -1), "0"},
		{12.4, "12.4"},
		{-3.1, "-3.1"},
		{1.5708, "1.5708"},
	}
	for _, test := range cases {
		encoded, err := CanonicalJSON(map[string]any{"v": test.value})
		if err != nil {
			t.Fatal(err)
		}
		if string(encoded) != `{"v":`+test.expected+`}` {
			t.Errorf("value %v = %s, want %s", test.value, encoded, test.expected)
		}
	}
}

func TestBindingDigestMatchesFixture(t *testing.T) {
	raw, err := os.ReadFile("../../../contracts/operations/v1/fixtures/valid/submit-navigate.json")
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]json.RawMessage
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatal(err)
	}
	var expected string
	if err := json.Unmarshal(document["binding_digest"], &expected); err != nil {
		t.Fatal(err)
	}
	var request SubmitOperationRequest
	if err := json.Unmarshal(raw, &request); err != nil {
		t.Fatal(err)
	}
	digest, err := BindingDigest(request)
	if err != nil {
		t.Fatal(err)
	}
	if digest != expected {
		t.Fatalf("binding digest = %s, want %s", digest, expected)
	}
}

func TestFillDigests(t *testing.T) {
	request := SubmitOperationRequest{
		OperationType: "ros2.nav2.navigate_to_pose@1",
		Scope:         Scope{TenantID: "tenant"},
	}
	if err := FillDigests(&request); err != nil {
		t.Fatal(err)
	}
	if len(request.PayloadDigest) != len("sha256:")+64 ||
		len(request.BindingDigest) != len("sha256:")+64 {
		t.Fatalf("digests = %q %q", request.PayloadDigest, request.BindingDigest)
	}
	if request.PayloadDigest[:7] != "sha256:" || request.BindingDigest[:7] != "sha256:" {
		t.Fatalf("digests = %q %q", request.PayloadDigest, request.BindingDigest)
	}
}
