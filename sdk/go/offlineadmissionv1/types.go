// Package offlineadmissionv1 contains types for the proposed AEGIVELA
// offline-admission transport. The contract is not reviewed or enabled; this
// package provides encoding only and performs no signature or trust checks.
package offlineadmissionv1

import (
	"encoding/json"
	"time"
)

const (
	ContractVersion            = "v1.0"
	WireVersion                = "aegivela.io/offline-admission/v1.0"
	AegivelaRevision           = "684886a1f62e5e8f2159250a584ccd8e4187f667"
	CapabilityOfflineAdmission = "offline_admission_v1"
	ProtocolVersion            = "akrovela.satellite/v1"
)

type Envelope struct {
	ProtocolVersion  string          `json:"protocol_version"`
	MessageType      string          `json:"message_type"`
	MessageID        string          `json:"message_id"`
	SessionID        string          `json:"session_id"`
	Sequence         int64           `json:"sequence"`
	TenantID         string          `json:"tenant_id"`
	WorkspaceID      string          `json:"workspace_id"`
	SiteID           string          `json:"site_id"`
	SatelliteID      string          `json:"satellite_id"`
	SourceEpoch      int64           `json:"source_epoch"`
	CoordinatorID    string          `json:"coordinator_id"`
	CoordinatorEpoch int64           `json:"coordinator_epoch"`
	SentAt           time.Time       `json:"sent_at"`
	Body             json.RawMessage `json:"body"`
}

type CapabilityAck struct {
	ProtocolVersion string   `json:"protocol_version"`
	Accepted        bool     `json:"accepted"`
	Capabilities    []string `json:"capabilities"`
}

type OfflineGrant struct {
	AegivelaRevision   string `json:"aegivela_revision"`
	SignedOfflineGrant string `json:"signed_offline_grant"`
	SignedParentGrant  string `json:"signed_parent_grant"`
	ArtifactDigest     string `json:"artifact_digest"`
}

type OfflineGrantAck struct {
	ArtifactDigest string    `json:"artifact_digest"`
	Status         string    `json:"status"`
	RecordedAt     time.Time `json:"recorded_at"`
	Reason         string    `json:"reason,omitempty"`
}

type ReservedResources struct {
	Operations     int64 `json:"operations"`
	DurationMS     int64 `json:"duration_ms"`
	GoalDistanceMM int64 `json:"goal_distance_mm"`
	PathLengthMM   int64 `json:"path_length_mm"`
}

type AdmissionFact struct {
	FactID                 string            `json:"fact_id"`
	AdmissionID            string            `json:"admission_id"`
	Namespace              string            `json:"namespace"`
	OfflineGrantID         string            `json:"offline_grant_id"`
	ArtifactDigest         string            `json:"artifact_digest"`
	OperationKey           string            `json:"operation_key"`
	AttemptID              string            `json:"attempt_id"`
	AuthorityBindingDigest string            `json:"authority_binding_digest"`
	AuthorityPayloadDigest string            `json:"authority_payload_digest"`
	OperationBindingDigest string            `json:"operation_binding_digest"`
	OperationPayloadDigest string            `json:"operation_payload_digest"`
	Sequence               int64             `json:"sequence"`
	State                  string            `json:"state"`
	ReservedResources      ReservedResources `json:"reserved_resources"`
	RecordedAt             time.Time         `json:"recorded_at"`
}

type AdmissionBatch struct {
	AegivelaRevision     string          `json:"aegivela_revision"`
	Namespace            string          `json:"namespace"`
	AdmissionSourceEpoch int64           `json:"admission_source_epoch"`
	Facts                []AdmissionFact `json:"facts"`
}

type AdmissionAck struct {
	Namespace            string `json:"namespace"`
	AdmissionSourceEpoch int64  `json:"admission_source_epoch"`
	AckedThroughSequence int64  `json:"acked_through_sequence"`
}

var MessageTypes = map[string]bool{
	"hello_ack": true, "offline_grant": true, "offline_grant_ack": true,
	"local_admission_batch": true, "local_admission_ack": true,
}
