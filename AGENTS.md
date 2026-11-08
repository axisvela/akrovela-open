# AKROVELA Open Rules

- Use installed Superpowers; see `superpowers/README.md`. Report unavailable tools.
- Only public contracts, SDKs, examples, public usage and approved binaries here.
  Internal product architecture/design/plans belong in EE.
- SDK code is Apache-2.0; never import AGPL core or cross-product internals.
- From `sdk/go/`: `go test ./...`, `go vet ./...`, `go build ./...`.
  From root: `python -m unittest discover -s tests -v`.
- Version public contracts before dependent changes; no nonexistent domain APIs.
- No credentials, tenant data or training data in examples or binary bundles.
- X.Y.Z only; align core/open. Root release tags and nested SDK module tags are
  different publication artifacts; neither is created automatically.
- Background any service tests. Do not commit, tag, push or publish unless asked.
