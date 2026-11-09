package operationsv1

import (
	"os"
	"path/filepath"
	"testing"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "contracts/operations/v1/fixtures/valid", name))
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestValidFixturesDecode(t *testing.T) {
	for _, name := range []string{"submit-navigate.json", "submit-small-floats.json"} {
		t.Run(name, func(t *testing.T) {
			value, err := Decode[SubmitOperationRequest](fixture(t, name))
			if err != nil || value.OperationKey == "" || value.Target.NodeID == "" {
				t.Fatalf("value=%+v err=%v", value, err)
			}
		})
	}
	if value, err := Decode[CancelRequest](fixture(t, "cancel-request.json")); err != nil || value.CancelRequestID == "" {
		t.Fatalf("cancel request err=%v", err)
	}
	if value, err := Decode[CancelResponse](fixture(t, "cancel-response.json")); err != nil || value.CancelRequestID == "" {
		t.Fatalf("cancel response err=%v", err)
	}
	if value, err := Decode[OperationView](fixture(t, "operation-view.json")); err != nil || value.OperationKey == "" {
		t.Fatalf("operation view err=%v", err)
	}
	for _, name := range []string{"receipt.json", "receipt-admission-rejected.json"} {
		if value, err := Decode[Receipt](fixture(t, name)); err != nil || value.ReceiptID == "" {
			t.Fatalf("%s err=%v", name, err)
		}
	}
}

func TestDecodeRejectsUnknownAndTrailing(t *testing.T) {
	raw := fixture(t, "cancel-request.json")
	withUnknown := append([]byte(nil), raw[:len(raw)-1]...)
	withUnknown = append(withUnknown, []byte(`,"extra":true}`)...)
	if _, err := Decode[CancelRequest](withUnknown); err == nil {
		t.Fatal("accepted unknown field")
	}
	if _, err := Decode[CancelRequest](append(append([]byte(nil), raw...), []byte("{}")...)); err == nil {
		t.Fatal("accepted trailing data")
	}
}
