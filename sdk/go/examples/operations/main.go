// Command operations exercises the reviewed operation client against a local
// AKROVELA core. It does not execute a robot command by itself.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	operationsv1 "github.com/axisvela/akrovela-open/sdk/go/operationsv1"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, out io.Writer) error {
	if len(args) == 0 {
		return errors.New(usage)
	}
	command := args[0]
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	baseURL := flags.String("base-url", "http://127.0.0.1:2046", "AKROVELA API base URL")
	token := flags.String("token", os.Getenv("AKROVELA_API_TOKEN"), "bearer token")
	operation := flags.String("operation", "", "operation key")
	requestPath := flags.String("request", "", "submit request JSON file, or - for stdin")
	receiptID := flags.String("receipt", "", "receipt id")
	after := flags.Int64("after", 0, "receipt after_sequence")
	cancelID := flags.String("cancel-id", "", "cancel request id")
	reason := flags.String("reason", "", "cancel reason code")
	deadline := flags.Int("deadline-ms", 500, "cancel deadline in milliseconds")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	client := operationsv1.NewClient(*baseURL, *token)
	ctx := context.Background()
	switch command {
	case "submit":
		request, err := readSubmitRequest(*requestPath, stdin)
		if err != nil {
			return err
		}
		if err := operationsv1.FillDigests(&request); err != nil {
			return err
		}
		view, err := client.Submit(ctx, request)
		if err != nil {
			return err
		}
		return writeJSON(out, view)
	case "get":
		if *operation == "" {
			return errors.New("-operation is required")
		}
		view, err := client.Get(ctx, *operation)
		if err != nil {
			return err
		}
		return writeJSON(out, view)
	case "cancel":
		if *operation == "" || *cancelID == "" {
			return errors.New("-operation and -cancel-id are required")
		}
		response, err := client.Cancel(ctx, *operation, operationsv1.CancelRequest{
			CancelRequestID: *cancelID, ReasonCode: *reason,
			ClientRequestedAt: time.Now().UTC(), DeadlineMS: *deadline,
		})
		if err != nil {
			return err
		}
		return writeJSON(out, response)
	case "receipts":
		if *operation == "" {
			return errors.New("-operation is required")
		}
		list, err := client.Receipts(ctx, *operation, *after)
		if err != nil {
			return err
		}
		return writeJSON(out, list)
	case "receipt":
		if *receiptID == "" {
			return errors.New("-receipt is required")
		}
		receipt, err := client.GetReceipt(ctx, *receiptID)
		if err != nil {
			return err
		}
		return writeJSON(out, receipt)
	default:
		return fmt.Errorf("unknown command %q\n%s", command, usage)
	}
}

func readSubmitRequest(path string, stdin io.Reader) (operationsv1.SubmitOperationRequest, error) {
	if path == "" {
		return operationsv1.SubmitOperationRequest{}, errors.New("-request is required")
	}
	var raw []byte
	var err error
	if path == "-" {
		raw, err = io.ReadAll(io.LimitReader(stdin, 4<<20))
	} else {
		raw, err = os.ReadFile(path)
	}
	if err != nil {
		return operationsv1.SubmitOperationRequest{}, err
	}
	var request operationsv1.SubmitOperationRequest
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return operationsv1.SubmitOperationRequest{}, err
	}
	return request, nil
}

func writeJSON(out io.Writer, value any) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(encoded))
	return err
}

const usage = `usage:
  operations submit -request request.json [-base-url URL] [-token TOKEN]
  operations get -operation KEY [-base-url URL] [-token TOKEN]
  operations cancel -operation KEY -cancel-id ID [-reason CODE] [-deadline-ms MS]
  operations receipts -operation KEY [-after SEQUENCE]
  operations receipt -receipt ID`
