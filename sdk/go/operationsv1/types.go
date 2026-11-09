package operationsv1

import "time"

const Version = "v1"
const ProtocolVersion = "akrovela.local-operation/v1"

type Scope struct {
	TenantID    string `json:"tenant_id"`
	WorkspaceID string `json:"workspace_id"`
	SiteID      string `json:"site_id"`
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

type ObservationSubject struct {
	RobotRef       string `json:"robot_ref"`
	OperationKey   string `json:"operation_key"`
	AttemptID      string `json:"attempt_id"`
	ExecutionEpoch int64  `json:"execution_epoch"`
	MapRef         string `json:"map_ref"`
	TargetDigest   string `json:"target_digest"`
}

type ObservationReference struct {
	ObservationID     string              `json:"observation_id"`
	SourceID          string              `json:"source_id"`
	TrustDomain       string              `json:"trust_domain"`
	StateLabel        string              `json:"state_label"`
	ClockDomain       string              `json:"clock_domain"`
	SourceEpoch       int64               `json:"source_epoch,omitempty"`
	Seq               int64               `json:"seq"`
	ObservedAt        time.Time           `json:"observed_at"`
	ReceivedAt        time.Time           `json:"received_at"`
	MaxAgeMS          int64               `json:"max_age_ms,omitempty"`
	MaxClockSkewMS    int                 `json:"max_clock_skew_ms,omitempty"`
	FreshnessDecision string              `json:"freshness_decision"`
	Subject           *ObservationSubject `json:"subject,omitempty"`
	Digest            string              `json:"digest"`
}

type Goal struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Yaw     float64 `json:"yaw"`
	FrameID string  `json:"frame_id"`
}

type NavigateToPosePayload struct {
	Goal               Goal    `json:"goal"`
	PositionToleranceM float64 `json:"position_tolerance_m"`
	YawToleranceRad    float64 `json:"yaw_tolerance_rad"`
	BehaviorTree       string  `json:"behavior_tree"`
	PlannerID          string  `json:"planner_id"`
	ControllerID       string  `json:"controller_id"`
}

type SubmitOperationRequest struct {
	ProtocolVersion             string                 `json:"protocol_version"`
	OperationKey                string                 `json:"operation_key"`
	OperationType               string                 `json:"operation_type"`
	PayloadDigest               string                 `json:"payload_digest"`
	BindingDigest               string                 `json:"binding_digest"`
	Payload                     NavigateToPosePayload  `json:"payload"`
	Scope                       Scope                  `json:"scope"`
	Target                      Target                 `json:"target"`
	MissionRef                  string                 `json:"mission_ref"`
	AssignmentRef               string                 `json:"assignment_ref"`
	LeaseRef                    string                 `json:"lease_ref"`
	FencingToken                int64                  `json:"fencing_token"`
	SafetyEnvelopeRef           string                 `json:"safety_envelope_ref"`
	AuthorityRef                string                 `json:"authority_ref"`
	ApprovalRef                 string                 `json:"approval_ref"`
	Validity                    Validity               `json:"validity"`
	ResourceLimits              ResourceLimits         `json:"resource_limits"`
	ObservationRefs             []ObservationReference `json:"observation_refs"`
	IdempotencyRetentionSeconds int                    `json:"idempotency_retention_seconds"`
}

type AttemptError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}

type AttemptView struct {
	AttemptID      string        `json:"attempt_id"`
	NodeBootEpoch  int64         `json:"node_boot_epoch"`
	ExecutionEpoch int64         `json:"execution_epoch"`
	FencingToken   int64         `json:"fencing_token"`
	Status         string        `json:"status"`
	StartedAt      *time.Time    `json:"started_at,omitempty"`
	CompletedAt    *time.Time    `json:"completed_at,omitempty"`
	Error          *AttemptError `json:"error,omitempty"`
}

type ExecutorFacts struct {
	Adapter          string     `json:"adapter"`
	GoalID           string     `json:"goal_id"`
	TerminalState    string     `json:"terminal_state"`
	ActionTerminated bool       `json:"action_terminated"`
	StartedAt        *time.Time `json:"started_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
}

type VerifierResult struct {
	VerifierID            string    `json:"verifier_id"`
	Result                string    `json:"result"`
	EvaluatedAt           time.Time `json:"evaluated_at"`
	IndependentOfExecutor bool      `json:"independent_of_executor"`
	TrustDomain           string    `json:"trust_domain"`
	SharedFailureModes    []string  `json:"shared_failure_modes"`
	EvidenceRef           string    `json:"evidence_ref"`
}

type Receipt struct {
	ReceiptID          string                 `json:"receipt_id"`
	ReceiptSequence    int64                  `json:"receipt_sequence"`
	ReceiptType        string                 `json:"receipt_type"`
	OperationKey       string                 `json:"operation_key"`
	AttemptID          string                 `json:"attempt_id,omitempty"`
	EmittedAt          time.Time              `json:"emitted_at"`
	PayloadDigest      string                 `json:"payload_digest"`
	BindingDigest      string                 `json:"binding_digest"`
	FencingToken       int64                  `json:"fencing_token"`
	ContractMode       string                 `json:"contract_mode"`
	AdmissionStatus    string                 `json:"admission_status"`
	ExecutionStatus    string                 `json:"execution_status"`
	CancellationStatus string                 `json:"cancellation_status"`
	EffectStatus       string                 `json:"effect_status"`
	ExecutorFacts      *ExecutorFacts         `json:"executor_facts,omitempty"`
	Observations       []ObservationReference `json:"observations,omitempty"`
	Verifier           *VerifierResult        `json:"verifier,omitempty"`
	ConfirmationSource string                 `json:"confirmation_source,omitempty"`
	ConfirmationRef    string                 `json:"confirmation_ref,omitempty"`
	EvidenceRefs       []string               `json:"evidence_refs"`
	Error              *AttemptError          `json:"error,omitempty"`
	ReceiptDigest      string                 `json:"receipt_digest"`
}

type ReceiptList struct {
	Receipts          []Receipt `json:"receipts"`
	NextAfterSequence int64     `json:"next_after_sequence"`
}

type OperationView struct {
	ProtocolVersion             string                 `json:"protocol_version"`
	OperationKey                string                 `json:"operation_key"`
	OperationType               string                 `json:"operation_type"`
	PayloadDigest               string                 `json:"payload_digest"`
	BindingDigest               string                 `json:"binding_digest"`
	Payload                     NavigateToPosePayload  `json:"payload"`
	ObservationRefs             []ObservationReference `json:"observation_refs"`
	IdempotencyRetentionSeconds int                    `json:"idempotency_retention_seconds"`
	Scope                       Scope                  `json:"scope"`
	Target                      Target                 `json:"target"`
	MissionRef                  string                 `json:"mission_ref"`
	AssignmentRef               string                 `json:"assignment_ref"`
	LeaseRef                    string                 `json:"lease_ref"`
	FencingToken                int64                  `json:"fencing_token"`
	SafetyEnvelopeRef           string                 `json:"safety_envelope_ref"`
	AuthorityRef                string                 `json:"authority_ref"`
	ApprovalRef                 string                 `json:"approval_ref"`
	Validity                    Validity               `json:"validity"`
	ResourceLimits              ResourceLimits         `json:"resource_limits"`
	AdmissionStatus             string                 `json:"admission_status"`
	ExecutionStatus             string                 `json:"execution_status"`
	CancellationStatus          string                 `json:"cancellation_status"`
	EffectStatus                string                 `json:"effect_status"`
	ConfirmationSource          string                 `json:"confirmation_source,omitempty"`
	ConfirmationRef             string                 `json:"confirmation_ref,omitempty"`
	Attempts                    []AttemptView          `json:"attempts"`
	LatestReceipt               *Receipt               `json:"latest_receipt,omitempty"`
	ContractMode                string                 `json:"contract_mode"`
	CreatedAt                   time.Time              `json:"created_at"`
	UpdatedAt                   time.Time              `json:"updated_at"`
}

type CancelRequest struct {
	CancelRequestID   string    `json:"cancel_request_id"`
	ReasonCode        string    `json:"reason_code"`
	ClientRequestedAt time.Time `json:"client_requested_at"`
	DeadlineMS        int       `json:"deadline_ms"`
}

type CancelResponse struct {
	CancelRequestID    string    `json:"cancel_request_id"`
	CancelStatus       string    `json:"cancel_status"`
	ReceivedAt         time.Time `json:"received_at"`
	StopConfirmed      bool      `json:"stop_confirmed"`
	ConfirmationSource string    `json:"confirmation_source"`
	ReceiptID          string    `json:"receipt_id,omitempty"`
}

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}

type Error struct {
	Error ErrorDetail `json:"error"`
}
