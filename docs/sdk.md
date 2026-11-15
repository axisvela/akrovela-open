# Public Go SDK

The Go SDK lives under `sdk/go` as module
`github.com/axisvela/akrovela-open/sdk/go`. It is dependency-free and does not
import AKROVELA core, EE or other product internals.

| Package | Contract | Stability |
| --- | --- | --- |
| `akrovela` (root) | `/version` client | Reviewed, minimal |
| `operationsv1` | `contracts/operations/v1/operations.openapi.json` | Reviewed types, strict JSON codec, JCS digest helpers and HTTP client |
| `satellitev1` | `contracts/satellite/v1/satellite.schema.json` | Reviewed types and strict JSON codec |
| `offlineadmissionv1` | `contracts/satellite/v1/offline-admission.schema.json` | **Experimental**: contract status is proposed |

## What the SDK Does

- Provides Go types mirroring the public JSON fields.
- `Decode`, `DecodeEnvelope` and `DecodeBody[T]` reject unknown fields, trailing
  data and wrong protocol/message identifiers.
- `Encode`/`EncodeEnvelope` marshal the same types.
- `operationsv1.Client` calls the reviewed operation submit/query/cancel/receipt
  endpoints with bearer authentication and typed error mapping.
- `PayloadDigest`, `BindingDigest` and `FillDigests` compute the RFC 8785 JCS
  SHA-256 digests over the declared operation payload and binding fields.

## What It Does Not Do

- It is not an authorization or policy implementation.
- It does not verify offline-admission signatures, trust bundles, parent lineage,
  reservations, revocation or device trust. Use the reviewed AKROVELA consumer
  path for those decisions.
- It does not synthesize Mission/lease/fence values, retry effects or treat a
  cancel response as a stop confirmation. Digest helpers cover only the
  documented operation payload and binding objects.
- `offlineadmissionv1` must not be treated as a stable or enabled contract.

## Validation Limits

Strict decoding proves structure and unknown-field rejection only. Required-field
presence, patterns, ranges, conditional rules and cross-message semantics remain
subject to the contract owner and consumer validation. Date-time fields use Go
`time.Time` and therefore accept RFC 3339.

## Example

```go
import (
    "context"
    "os"
    "time"

    operationsv1 "github.com/axisvela/akrovela-open/sdk/go/operationsv1"
    satellitev1 "github.com/axisvela/akrovela-open/sdk/go/satellitev1"
)

raw, _ := os.ReadFile("submit.json")
envelope, err := satellitev1.DecodeEnvelope(raw)
if err != nil {
    return err
}
submit, err := satellitev1.DecodeBody[satellitev1.SubmitBody](envelope)

client := operationsv1.NewClient("http://127.0.0.1:2046", "local-token")
view, err := client.Get(context.Background(), "op:example")

request := operationsv1.SubmitOperationRequest{
    OperationType: "ros2.nav2.navigate_to_pose@1",
    // ... fill the reviewed operation fields ...
}
if err := operationsv1.FillDigests(&request); err != nil {
    return err
}
submitted, err := client.Submit(context.Background(), request)
final, err := client.Wait(context.Background(), submitted.OperationKey, 500*time.Millisecond)
```

## Verification

From `sdk/go`: `go test ./...`, `go vet ./...`, `go build ./...`.
From the repository root: `python -m unittest discover -s tests -v`.
`sdk/go/sdkparity` compares schema properties with Go json tags in both directions
so an unreviewed field cannot silently appear in the SDK. The runnable
`examples/operations` command submits, queries, cancels and lists receipts
through the public client.

## Versioning and Compatibility

The nested module uses `sdk/go/vX.Y.Z` tags, separate from product root tags.
While the product is `0.x`, SDK breaking changes are allowed in a minor release,
but must be documented in `CHANGELOG.md` and reflected in fixture tests.

- `operationsv1` and `satellitev1` follow reviewed public contracts.
- `offlineadmissionv1` is experimental and may change without a contract review of
  its own; do not publish it as stable.
- A module tag is a publication decision, not an automatic consequence of a
  product commit. It requires a clean-checkout `go test/vet/build` run and release
  notes that exclude production/device/enablement claims.

No tag or registry availability is currently claimed.



