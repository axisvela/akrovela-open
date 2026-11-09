# Public Go SDK

The Go SDK lives under `sdk/go` as module
`github.com/axisvela/akrovela-open/sdk/go`. It is dependency-free and does not
import AKROVELA core, EE or other product internals.

| Package | Contract | Stability |
| --- | --- | --- |
| `akrovela` (root) | `/version` client | Reviewed, minimal |
| `operationsv1` | `contracts/operations/v1/operations.openapi.json` | Reviewed types and strict JSON codec |
| `satellitev1` | `contracts/satellite/v1/satellite.schema.json` | Reviewed types and strict JSON codec |
| `offlineadmissionv1` | `contracts/satellite/v1/offline-admission.schema.json` | **Experimental**: contract status is proposed |

## What the SDK Does

- Provides Go types mirroring the public JSON fields.
- `Decode`, `DecodeEnvelope` and `DecodeBody[T]` reject unknown fields, trailing
  data and wrong protocol/message identifiers.
- `Encode`/`EncodeEnvelope` marshal the same types.

## What It Does Not Do

- It is not an authorization or policy implementation.
- It does not verify offline-admission signatures, trust bundles, parent lineage,
  reservations, revocation or device trust. Use the reviewed AKROVELA consumer
  path for those decisions.
- It does not compute binding or payload digests, synthesize Mission/lease/fence
  values, or retry effects.
- `offlineadmissionv1` must not be treated as a stable or enabled contract.

## Validation Limits

Strict decoding proves structure and unknown-field rejection only. Required-field
presence, patterns, ranges, conditional rules and cross-message semantics remain
subject to the contract owner and consumer validation. Date-time fields use Go
`time.Time` and therefore accept RFC 3339.

## Example

```go
import (
    "os"

    satellitev1 "github.com/axisvela/akrovela-open/sdk/go/satellitev1"
)

raw, _ := os.ReadFile("submit.json")
envelope, err := satellitev1.DecodeEnvelope(raw)
if err != nil {
    return err
}
submit, err := satellitev1.DecodeBody[satellitev1.SubmitBody](envelope)
```

## Verification

From `sdk/go`: `go test ./...`, `go vet ./...`, `go build ./...`.
From the repository root: `python -m unittest discover -s tests -v`.

Publishing the nested module needs `sdk/go/vX.Y.Z`, separately from product root
tags. No tag or registry availability is claimed here.
