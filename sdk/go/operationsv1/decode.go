package operationsv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var ErrInvalidDocument = errors.New("operationsv1: invalid document")

func Decode[T any](encoded []byte) (T, error) {
	var value T
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return value, ErrInvalidDocument
	}
	if _, err := decoder.Token(); err != io.EOF {
		return value, ErrInvalidDocument
	}
	return value, nil
}

func Encode(value any) ([]byte, error) {
	return json.Marshal(value)
}
