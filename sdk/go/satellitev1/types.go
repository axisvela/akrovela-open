package satellitev1

import (
	"encoding/json"
	"time"
)

const (
	ProtocolVersion = "akrovela.satellite/v1"
	ContractVersion = "v1"
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

type Target struct {
	NodeID         string `json:"node_id"`
	NodeBootEpoch  int64  `json:"node_boot_epoch"`
	RobotRef       string `json:"robot_ref"`
	Adapter        string `json:"adapter"`
	AdapterVersion string `json:"adapter_version"`
	MapRef         string `json:"map_ref"`
	FrameID        string `json:"frame_id"`
}

type Validity struct {
	NotBefore           time.Time `json:"not_before"`
	ExpiresAt           time.Time `json:"expires_at"`
	ClockDomain         string    `json:"clock_domain"`
	MaxClockSkewMS      int       `json:"max_clock_skew_ms"`
	RevalidateAtRelease bool      `json:"revalidate_at_release"`
}

type ResourceLimits struct {
	MaxDurationMS    int64   `json:"max_duration_ms"`
	MaxGoalDistanceM float64 `json:"max_goal_distance_m"`
	MaxPathLengthM   float64 `json:"max_path_length_m"`
	MaxSpeedMPS      float64 `json:"max_speed_mps"`
	MaxAttempts      int     `json:"max_attempts"`
	MaxRetries       int     `json:"max_retries"`
	MaxPayloadBytes  int     `json:"max_payload_bytes"`
}

type Goal struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Yaw     float64 `json:"yaw"`
	FrameID string  `json:"frame_id"`
}

type NavigatePayload struct {
	Goal               Goal    `json:"goal"`
	PositionToleranceM float64 `json:"position_tolerance_m"`
	YawToleranceRad    float64 `json:"yaw_tolerance_rad"`
	BehaviorTree       string  `json:"behavior_tree"`
	PlannerID          string  `json:"planner_id"`
	ControllerID       string  `json:"controller_id"`
}

type SubmitBody struct {
	OperationKey                string          `json:"operation_key"`
	AttemptID                   string          `json:"attempt_id"`
	OperationType               string          `json:"operation_type"`
	PayloadDigest               string          `json:"payload_digest"`
	BindingDigest               string          `json:"binding_digest"`
	Payload                     NavigatePayload `json:"payload"`
	Target                      Target          `json:"target"`
	MissionRef                  string          `json:"mission_ref"`
	AssignmentRef               string          `json:"assignment_ref"`
	LeaseRef                    string          `json:"lease_ref"`
	FencingToken                int64           `json:"fencing_token"`
	SafetyEnvelopeRef           string          `json:"safety_envelope_ref"`
	AuthorityRef                string          `json:"authority_ref"`
	ApprovalRef                 string          `json:"approval_ref"`
	Validity                    Validity        `json:"validity"`
	ResourceLimits              ResourceLimits  `json:"resource_limits"`
	ObservationRefs             []any           `json:"observation_refs"`
	IdempotencyRetentionSeconds int             `json:"idempotency_retention_seconds"`
	DispatchNotAfter            time.Time       `json:"dispatch_not_after"`
	ContractMode                string          `json:"contract_mode"`
}

type SubmitResultBody struct {
	OperationKey string    `json:"operation_key"`
	AttemptID    string    `json:"attempt_id"`
	CommandID    string    `json:"command_id"`
	Status       string    `json:"status"`
	Reason       string    `json:"reason,omitempty"`
	AcceptedAt   time.Time `json:"accepted_at"`
	SourceEpoch  int64     `json:"source_epoch"`
}

type StatusBody struct {
	OperationKey string `json:"operation_key"`
	AttemptID    string `json:"attempt_id"`
}

type StatusResultBody struct {
	OperationKey string    `json:"operation_key"`
	AttemptID    string    `json:"attempt_id"`
	State        string    `json:"state"`
	Reason       string    `json:"reason,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type HelloBody struct {
	ProtocolMin  string   `json:"protocol_min"`
	ProtocolMax  string   `json:"protocol_max"`
	Capabilities []string `json:"capabilities"`
}

type HelloAckBody struct {
	ProtocolVersion string `json:"protocol_version"`
	Accepted        bool   `json:"accepted"`
	Reason          string `json:"reason,omitempty"`
}

type ResumeBody struct {
	Token                string `json:"token"`
	LastUplinkSequence   int64  `json:"last_uplink_sequence"`
	LastDownlinkSequence int64  `json:"last_downlink_sequence"`
	SourceEpoch          int64  `json:"source_epoch"`
	CoordinatorEpoch     int64  `json:"coordinator_epoch"`
}

type ResumeResultBody struct {
	Accepted          bool   `json:"accepted"`
	DrainState        string `json:"drain_state"`
	UplinkHighWater   int64  `json:"uplink_high_water"`
	DownlinkHighWater int64  `json:"downlink_high_water"`
	Reason            string `json:"reason,omitempty"`
}

type CancelBody struct {
	OperationKey      string    `json:"operation_key"`
	AttemptID         string    `json:"attempt_id"`
	CancelRequestID   string    `json:"cancel_request_id"`
	ReasonCode        string    `json:"reason_code"`
	ClientRequestedAt time.Time `json:"client_requested_at"`
	DeadlineMS        int       `json:"deadline_ms"`
}

type CancelResultBody struct {
	CancelRequestID    string    `json:"cancel_request_id"`
	Status             string    `json:"status"`
	ConfirmationSource string    `json:"confirmation_source,omitempty"`
	ConfirmationRef    string    `json:"confirmation_ref,omitempty"`
	AssertedAt         time.Time `json:"asserted_at"`
}

type SatelliteReceiptBody struct {
	SatelliteReceiptID string    `json:"satellite_receipt_id"`
	OperationKey       string    `json:"operation_key"`
	AttemptID          string    `json:"attempt_id"`
	ReceiptType        string    `json:"receipt_type"`
	Sequence           int64     `json:"sequence"`
	Digest             string    `json:"digest"`
	EmittedAt          time.Time `json:"emitted_at"`
}

type ReceiptBatchBody struct {
	SatelliteID string                 `json:"satellite_id"`
	Receipts    []SatelliteReceiptBody `json:"receipts"`
}

type ReceiptAckBody struct {
	SatelliteID          string `json:"satellite_id"`
	AckedThroughSequence int64  `json:"acked_through_sequence"`
}

type HeartbeatBody struct {
	State      string `json:"state"`
	QueueDepth int    `json:"queue_depth"`
}

type ResourceBudget struct {
	MaxDurationMS    int64   `json:"max_duration_ms"`
	MaxGoalDistanceM float64 `json:"max_goal_distance_m"`
	MaxPathLengthM   float64 `json:"max_path_length_m"`
	MaxSpeedMPS      float64 `json:"max_speed_mps"`
	MaxPayloadBytes  int     `json:"max_payload_bytes"`
}

type PreReleaseBody struct {
	SubmitBody
	AllocationID   string         `json:"allocation_id"`
	CountBudget    int            `json:"count_budget"`
	ResourceBudget ResourceBudget `json:"resource_budget"`
}

type PreReleaseAckBody struct {
	AllocationID string    `json:"allocation_id"`
	Accepted     bool      `json:"accepted"`
	Reason       string    `json:"reason,omitempty"`
	PersistedAt  time.Time `json:"persisted_at"`
}

type ConsumptionBody struct {
	ConsumptionID string    `json:"consumption_id"`
	AllocationID  string    `json:"allocation_id"`
	OperationKey  string    `json:"operation_key"`
	AttemptID     string    `json:"attempt_id"`
	ConsumedAt    time.Time `json:"consumed_at"`
	Sequence      int64     `json:"sequence"`
	Digest        string    `json:"digest"`
	ReceiptRef    string    `json:"receipt_ref,omitempty"`
}

type ConsumptionBatchBody struct {
	SatelliteID  string            `json:"satellite_id"`
	Consumptions []ConsumptionBody `json:"consumptions"`
}

type ConsumptionAckBody struct {
	SatelliteID          string `json:"satellite_id"`
	AckedThroughSequence int64  `json:"acked_through_sequence"`
}

var MessageTypes = map[string]bool{
	"hello": true, "hello_ack": true, "resume": true, "resume_result": true,
	"submit": true, "submit_result": true, "status": true, "status_result": true,
	"cancel": true, "cancel_result": true, "receipt_batch": true, "receipt_ack": true,
	"pre_release": true, "pre_release_ack": true, "consumption_batch": true, "consumption_ack": true,
	"heartbeat": true,
}
