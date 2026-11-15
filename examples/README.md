# Examples

Run `go run ./examples/version http://127.0.0.1:2046` from `sdk/go/` after starting
the core. The runnable example is inside the SDK module so no unpublished external
module dependency or local `replace` directive is needed.

Run `go run ./examples/operations submit -request request.json` from `sdk/go/`
to fill the reviewed payload/binding digests and submit through the public
`operationsv1.Client`. The command also supports `get`, `cancel`, `receipts` and
`receipt`; pass `-base-url` and `-token` or set `AKROVELA_API_TOKEN`.
