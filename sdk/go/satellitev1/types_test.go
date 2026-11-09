package satellitev1

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func decodeFixture(t *testing.T, name string) Envelope {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "contracts/satellite/v1/fixtures/valid", name))
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]json.RawMessage
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatal(err)
	}
	delete(document, "_schema")
	encoded, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := DecodeEnvelope(encoded)
	if err != nil {
		t.Fatal(err)
	}
	return envelope
}

func TestValidFixturesDecode(t *testing.T) {
	for _, name := range []string{
		"hello.json", "hello-ack.json", "resume.json", "resume-result.json",
		"submit.json", "submit-result.json", "status.json", "status-result.json",
		"cancel.json", "cancel-result.json", "receipt-batch.json", "receipt-ack.json",
		"pre-release.json", "pre-release-ack.json", "consumption-batch.json", "consumption-ack.json",
		"heartbeat.json",
	} {
		t.Run(name, func(t *testing.T) {
			envelope := decodeFixture(t, name)
			if envelope.MessageID == "" || envelope.SatelliteID == "" {
				t.Fatalf("envelope=%+v", envelope)
			}
			switch envelope.MessageType {
			case "hello":
				value, err := DecodeBody[HelloBody](envelope)
				if err != nil || value.ProtocolMin == "" {
					t.Fatalf("hello err=%v", err)
				}
			case "hello_ack":
				if _, err := DecodeBody[HelloAckBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "resume":
				if value, err := DecodeBody[ResumeBody](envelope); err != nil || value.Token == "" {
					t.Fatalf("resume err=%v", err)
				}
			case "resume_result":
				if _, err := DecodeBody[ResumeResultBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "submit":
				if value, err := DecodeBody[SubmitBody](envelope); err != nil || value.OperationKey == "" || value.Target.NodeID == "" {
					t.Fatalf("submit err=%v", err)
				}
			case "submit_result":
				if value, err := DecodeBody[SubmitResultBody](envelope); err != nil || value.CommandID == "" {
					t.Fatalf("submit result err=%v", err)
				}
			case "status":
				if _, err := DecodeBody[StatusBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "status_result":
				if _, err := DecodeBody[StatusResultBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "cancel":
				if value, err := DecodeBody[CancelBody](envelope); err != nil || value.CancelRequestID == "" {
					t.Fatalf("cancel err=%v", err)
				}
			case "cancel_result":
				if _, err := DecodeBody[CancelResultBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "receipt_batch":
				if value, err := DecodeBody[ReceiptBatchBody](envelope); err != nil || len(value.Receipts) == 0 {
					t.Fatalf("receipt batch err=%v", err)
				}
			case "receipt_ack":
				if _, err := DecodeBody[ReceiptAckBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "pre_release":
				if value, err := DecodeBody[PreReleaseBody](envelope); err != nil || value.AllocationID == "" {
					t.Fatalf("pre-release err=%v", err)
				}
			case "pre_release_ack":
				if _, err := DecodeBody[PreReleaseAckBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "consumption_batch":
				if value, err := DecodeBody[ConsumptionBatchBody](envelope); err != nil || len(value.Consumptions) == 0 {
					t.Fatalf("consumption batch err=%v", err)
				}
			case "consumption_ack":
				if _, err := DecodeBody[ConsumptionAckBody](envelope); err != nil {
					t.Fatal(err)
				}
			case "heartbeat":
				if _, err := DecodeBody[HeartbeatBody](envelope); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestInvalidFixturesRejected(t *testing.T) {
	for _, name := range []string{"submit-unknown-field.json", "submit-wrong-version.json"} {
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("..", "..", "..", "contracts/satellite/v1/fixtures/invalid", name))
			if err != nil {
				t.Fatal(err)
			}
			var document map[string]json.RawMessage
			if err := json.Unmarshal(raw, &document); err != nil {
				t.Fatal(err)
			}
			delete(document, "_schema")
			encoded, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			envelope, err := DecodeEnvelope(encoded)
			if err != nil {
				return
			}
			if _, err := DecodeBody[SubmitBody](envelope); err == nil {
				t.Fatal("invalid submit body accepted")
			}
		})
	}
}
