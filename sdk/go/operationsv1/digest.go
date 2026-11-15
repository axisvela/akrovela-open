package operationsv1

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

func PayloadDigest(request SubmitOperationRequest) (string, error) {
	return digestValue(struct {
		OperationType string                `json:"operation_type"`
		Payload       NavigateToPosePayload `json:"payload"`
	}{OperationType: request.OperationType, Payload: request.Payload})
}

func BindingDigest(request SubmitOperationRequest) (string, error) {
	observationRefs := request.ObservationRefs
	if observationRefs == nil {
		observationRefs = []ObservationReference{}
	}
	return digestValue(struct {
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
	}{
		Scope: request.Scope, Target: request.Target, MissionRef: request.MissionRef,
		AssignmentRef: request.AssignmentRef, LeaseRef: request.LeaseRef,
		FencingToken: request.FencingToken, SafetyEnvelopeRef: request.SafetyEnvelopeRef,
		AuthorityRef: request.AuthorityRef, ApprovalRef: request.ApprovalRef,
		Validity: request.Validity, ResourceLimits: request.ResourceLimits,
		ObservationRefs: observationRefs, IdempotencyRetentionSeconds: request.IdempotencyRetentionSeconds,
	})
}

func FillDigests(request *SubmitOperationRequest) error {
	if request == nil {
		return fmt.Errorf("operationsv1: request is required")
	}
	payloadDigest, err := PayloadDigest(*request)
	if err != nil {
		return err
	}
	bindingDigest, err := BindingDigest(*request)
	if err != nil {
		return err
	}
	request.PayloadDigest = payloadDigest
	request.BindingDigest = bindingDigest
	return nil
}

func CanonicalJSON(value any) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	if err := writeCanonical(&buffer, decoded); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func digestValue(value any) (string, error) {
	canonical, err := CanonicalJSON(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func writeCanonical(buffer *bytes.Buffer, value any) error {
	switch typed := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		buffer.WriteByte('{')
		for index, key := range keys {
			if index > 0 {
				buffer.WriteByte(',')
			}
			writeString(buffer, key)
			buffer.WriteByte(':')
			if err := writeCanonical(buffer, typed[key]); err != nil {
				return err
			}
		}
		buffer.WriteByte('}')
	case []any:
		buffer.WriteByte('[')
		for index, item := range typed {
			if index > 0 {
				buffer.WriteByte(',')
			}
			if err := writeCanonical(buffer, item); err != nil {
				return err
			}
		}
		buffer.WriteByte(']')
	case string:
		writeString(buffer, typed)
	case json.Number:
		return writeNumber(buffer, typed.String())
	case bool:
		if typed {
			buffer.WriteString("true")
		} else {
			buffer.WriteString("false")
		}
	case nil:
		buffer.WriteString("null")
	default:
		return fmt.Errorf("operationsv1: unsupported canonical JSON value %T", value)
	}
	return nil
}

func writeString(buffer *bytes.Buffer, value string) {
	buffer.WriteByte('"')
	for _, character := range value {
		switch character {
		case '"':
			buffer.WriteString(`\"`)
		case '\\':
			buffer.WriteString(`\\`)
		case '\b':
			buffer.WriteString(`\b`)
		case '\f':
			buffer.WriteString(`\f`)
		case '\n':
			buffer.WriteString(`\n`)
		case '\r':
			buffer.WriteString(`\r`)
		case '\t':
			buffer.WriteString(`\t`)
		default:
			if character < 0x20 {
				fmt.Fprintf(buffer, `\u%04x`, character)
			} else {
				buffer.WriteRune(character)
			}
		}
	}
	buffer.WriteByte('"')
}

func writeNumber(buffer *bytes.Buffer, raw string) error {
	if raw == "0" || raw == "-0" {
		buffer.WriteByte('0')
		return nil
	}
	if !strings.ContainsAny(raw, ".eE") {
		if _, err := strconv.ParseInt(raw, 10, 64); err == nil {
			buffer.WriteString(raw)
			return nil
		}
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return err
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("operationsv1: non-finite number")
	}
	if value == 0 {
		buffer.WriteByte('0')
		return nil
	}
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	parts := strings.SplitN(strconv.FormatFloat(value, 'e', -1, 64), "e", 2)
	mantissa := parts[0]
	exponent, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	digits := strings.Replace(mantissa, ".", "", 1)
	fracDigits := 0
	if index := strings.IndexByte(mantissa, '.'); index >= 0 {
		fracDigits = len(mantissa) - index - 1
	}
	exponent -= fracDigits
	for len(digits) > 1 && strings.HasSuffix(digits, "0") {
		digits = digits[:len(digits)-1]
		exponent++
	}
	length := len(digits)
	point := exponent + length
	var out string
	switch {
	case length <= point && point <= 21:
		out = digits + strings.Repeat("0", point-length)
	case 0 < point && point <= 21:
		out = digits[:point] + "." + digits[point:]
	case -6 < point && point <= 0:
		out = "0." + strings.Repeat("0", -point) + digits
	default:
		out = digits
		if length > 1 {
			out = digits[:1] + "." + digits[1:]
		}
		exponentValue := point - 1
		if exponentValue >= 0 {
			out += "e+" + strconv.Itoa(exponentValue)
		} else {
			out += "e-" + strconv.Itoa(-exponentValue)
		}
	}
	buffer.WriteString(sign + out)
	return nil
}
