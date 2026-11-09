package satellitev1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var ErrInvalidDocument = errors.New("satellitev1: invalid document")

func DecodeEnvelope(encoded []byte) (Envelope, error) {
	var envelope Envelope
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return envelope, ErrInvalidDocument
	}
	if _, err := decoder.Token(); err != io.EOF {
		return envelope, ErrInvalidDocument
	}
	if envelope.ProtocolVersion != ProtocolVersion || !MessageTypes[envelope.MessageType] {
		return envelope, ErrInvalidDocument
	}
	return envelope, nil
}

func DecodeBody[T any](envelope Envelope) (T, error) {
	var body T
	decoder := json.NewDecoder(bytes.NewReader(envelope.Body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		return body, ErrInvalidDocument
	}
	if _, err := decoder.Token(); err != io.EOF {
		return body, ErrInvalidDocument
	}
	return body, nil
}

func EncodeEnvelope(envelope Envelope) ([]byte, error) {
	return json.Marshal(envelope)
}
