package sdkparity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	offlineadmissionv1 "github.com/axisvela/akrovela-open/sdk/go/offlineadmissionv1"
	operationsv1 "github.com/axisvela/akrovela-open/sdk/go/operationsv1"
	satellitev1 "github.com/axisvela/akrovela-open/sdk/go/satellitev1"
)

func schemaDefs(t *testing.T, path, section string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", path))
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatal(err)
	}
	defs, ok := document[section].(map[string]any)
	if !ok {
		t.Fatalf("%s: missing %s", path, section)
	}
	return defs
}

func jsonTags(t *testing.T, value any) map[string]bool {
	t.Helper()
	fields := map[string]bool{}
	collectTags(t, reflect.TypeOf(value), fields)
	return fields
}

func collectTags(t *testing.T, kind reflect.Type, fields map[string]bool) {
	t.Helper()
	for index := 0; index < kind.NumField(); index++ {
		field := kind.Field(index)
		if field.Anonymous {
			collectTags(t, field.Type, fields)
			continue
		}
		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		if tag == "" || tag == "-" {
			continue
		}
		fields[tag] = true
	}
}

func assertParity(t *testing.T, defs map[string]any, types map[string]any) {
	t.Helper()
	for name, value := range types {
		t.Run(name, func(t *testing.T) {
			definition, ok := defs[name].(map[string]any)
			if !ok {
				t.Fatalf("missing schema definition %s", name)
			}
			properties, _ := definition["properties"].(map[string]any)
			expected := map[string]bool{}
			for key := range properties {
				expected[key] = true
			}
			got := jsonTags(t, value)
			for key := range expected {
				if !got[key] {
					t.Errorf("schema field %q missing from Go type", key)
				}
			}
			for key := range got {
				if !expected[key] {
					t.Errorf("Go field %q missing from schema", key)
				}
			}
		})
	}
}

func TestOperationsParity(t *testing.T) {
	defs := schemaDefs(t, "contracts/operations/v1/operations.openapi.json", "components")
	schemas, ok := defs["schemas"].(map[string]any)
	if !ok {
		t.Fatal("missing components/schemas")
	}
	assertParity(t, schemas, map[string]any{
		"Scope":                  operationsv1.Scope{},
		"Target":                 operationsv1.Target{},
		"Validity":               operationsv1.Validity{},
		"ResourceLimits":         operationsv1.ResourceLimits{},
		"ObservationSubject":     operationsv1.ObservationSubject{},
		"ObservationReference":   operationsv1.ObservationReference{},
		"NavigateToPosePayload":  operationsv1.NavigateToPosePayload{},
		"SubmitOperationRequest": operationsv1.SubmitOperationRequest{},
		"AttemptView":            operationsv1.AttemptView{},
		"ExecutorFacts":          operationsv1.ExecutorFacts{},
		"VerifierResult":         operationsv1.VerifierResult{},
		"Receipt":                operationsv1.Receipt{},
		"ReceiptList":            operationsv1.ReceiptList{},
		"OperationView":          operationsv1.OperationView{},
		"CancelRequest":          operationsv1.CancelRequest{},
		"CancelResponse":         operationsv1.CancelResponse{},
		"Error":                  operationsv1.Error{},
	})
}

func TestSatelliteParity(t *testing.T) {
	defs := schemaDefs(t, "contracts/satellite/v1/satellite.schema.json", "$defs")
	assertParity(t, defs, map[string]any{
		"EnvelopeBase":         satellitev1.Envelope{},
		"Target":               satellitev1.Target{},
		"Validity":             satellitev1.Validity{},
		"ResourceLimits":       satellitev1.ResourceLimits{},
		"NavigatePayload":      satellitev1.NavigatePayload{},
		"SubmitBody":           satellitev1.SubmitBody{},
		"SubmitResultBody":     satellitev1.SubmitResultBody{},
		"CancelBody":           satellitev1.CancelBody{},
		"CancelResultBody":     satellitev1.CancelResultBody{},
		"ReceiptBatchBody":     satellitev1.ReceiptBatchBody{},
		"HeartbeatBody":        satellitev1.HeartbeatBody{},
		"HelloBody":            satellitev1.HelloBody{},
		"StatusBody":           satellitev1.StatusBody{},
		"StatusResultBody":     satellitev1.StatusResultBody{},
		"HelloAckBody":         satellitev1.HelloAckBody{},
		"ResumeBody":           satellitev1.ResumeBody{},
		"ResumeResultBody":     satellitev1.ResumeResultBody{},
		"ReceiptAckBody":       satellitev1.ReceiptAckBody{},
		"ResourceBudget":       satellitev1.ResourceBudget{},
		"PreReleaseBody":       satellitev1.PreReleaseBody{},
		"PreReleaseAckBody":    satellitev1.PreReleaseAckBody{},
		"ConsumptionBody":      satellitev1.ConsumptionBody{},
		"ConsumptionBatchBody": satellitev1.ConsumptionBatchBody{},
		"ConsumptionAckBody":   satellitev1.ConsumptionAckBody{},
	})
}

func TestOfflineAdmissionParity(t *testing.T) {
	defs := schemaDefs(t, "contracts/satellite/v1/offline-admission.schema.json", "$defs")
	assertParity(t, defs, map[string]any{
		"Envelope":          offlineadmissionv1.Envelope{},
		"CapabilityAck":     offlineadmissionv1.CapabilityAck{},
		"OfflineGrant":      offlineadmissionv1.OfflineGrant{},
		"OfflineGrantAck":   offlineadmissionv1.OfflineGrantAck{},
		"ReservedResources": offlineadmissionv1.ReservedResources{},
		"AdmissionFact":     offlineadmissionv1.AdmissionFact{},
		"AdmissionBatch":    offlineadmissionv1.AdmissionBatch{},
		"AdmissionAck":      offlineadmissionv1.AdmissionAck{},
	})
}
