package satellitev1_test

import (
	"encoding/json"
	"fmt"
	"time"

	satellitev1 "github.com/axisvela/akrovela-open/sdk/go/satellitev1"
)

func ExampleDecodeEnvelope() {
	body, _ := json.Marshal(satellitev1.HeartbeatBody{State: "ready", QueueDepth: 2})
	encoded, _ := json.Marshal(satellitev1.Envelope{
		ProtocolVersion: satellitev1.ProtocolVersion, MessageType: "heartbeat", MessageID: "msg:one",
		SessionID: "session:one", Sequence: 1, TenantID: "tenant-a", WorkspaceID: "wh3", SiteID: "site-wh3",
		SatelliteID: "node-1", SourceEpoch: 1, CoordinatorID: "coordinator-1", CoordinatorEpoch: 1,
		SentAt: time.Date(2026, 11, 9, 10, 0, 0, 0, time.UTC), Body: body,
	})
	envelope, err := satellitev1.DecodeEnvelope(encoded)
	if err != nil {
		panic(err)
	}
	heartbeat, err := satellitev1.DecodeBody[satellitev1.HeartbeatBody](envelope)
	if err != nil {
		panic(err)
	}
	fmt.Println(heartbeat.State)
	// Output: ready
}
