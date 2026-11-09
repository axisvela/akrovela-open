package offlineadmissionv1

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func fixture(t *testing.T, name string) Envelope {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "contracts/satellite/v1/offline-admission/fixtures/valid", name))
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := DecodeEnvelope(raw)
	if err != nil {
		t.Fatal(err)
	}
	return envelope
}

func TestValidFixturesDecode(t *testing.T) {
	if value, err := DecodeBody[CapabilityAck](fixture(t, "hello-ack.json")); err != nil || len(value.Capabilities) != 1 {
		t.Fatalf("capability ack err=%v", err)
	}
	if value, err := DecodeBody[OfflineGrant](fixture(t, "offline-grant.json")); err != nil || value.ArtifactDigest == "" {
		t.Fatalf("offline grant err=%v", err)
	}
	if _, err := DecodeBody[OfflineGrantAck](fixture(t, "offline-grant-ack.json")); err != nil {
		t.Fatal(err)
	}
	if value, err := DecodeBody[AdmissionBatch](fixture(t, "local-admission-batch.json")); err != nil || len(value.Facts) != 1 {
		t.Fatalf("admission batch err=%v", err)
	}
	if _, err := DecodeBody[AdmissionAck](fixture(t, "local-admission-ack.json")); err != nil {
		t.Fatal(err)
	}
}

func TestUnsignedGrantFixtureRejected(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "contracts/satellite/v1/offline-admission/fixtures/invalid/unsigned-grant.json"))
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := DecodeEnvelope(raw)
	if err != nil {
		return
	}
	if _, err := DecodeBody[OfflineGrant](envelope); err == nil {
		t.Fatal("unsigned grant accepted")
	}
}

func TestDecodeRejectsUnknownMessage(t *testing.T) {
	raw, _ := json.Marshal(map[string]any{
		"protocol_version": ProtocolVersion, "message_type": "offline_grant", "message_id": "msg:x",
		"session_id": "session:x", "sequence": 1, "tenant_id": "t", "workspace_id": "w", "site_id": "s",
		"satellite_id": "sat", "source_epoch": 1, "coordinator_id": "c", "coordinator_epoch": 1,
		"sent_at": "2026-11-09T10:00:00Z", "body": map[string]any{"extra": true},
	})
	envelope, err := DecodeEnvelope(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeBody[OfflineGrant](envelope); err == nil {
		t.Fatal("unknown body field accepted")
	}
}
