package operationsv1_test

import (
	"fmt"

	operationsv1 "github.com/axisvela/akrovela-open/sdk/go/operationsv1"
)

func ExampleDecode() {
	request, err := operationsv1.Decode[operationsv1.CancelRequest]([]byte(
		`{"cancel_request_id":"cancel:one","reason_code":"operator_request",` +
			`"client_requested_at":"2026-11-09T10:00:00Z","deadline_ms":2000}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(request.CancelRequestID)
	// Output: cancel:one
}
