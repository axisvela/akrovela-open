# AKROVELA Open

Public integration surface for AKROVELA, a governed edge/site execution platform.

| Repository | Scope | License |
| --- | --- | --- |
| [akrovela-open](https://github.com/axisvela/akrovela-open) | Public API, SDK, examples and approved binary distribution | Apache-2.0 |
| [akrovela](https://github.com/axisvela/akrovela) | Go/PostgreSQL core | AGPL-3.0-only |
| [akrovela-ee](https://github.com/axisvela/akrovela-ee) | Enterprise extensions and internal product documentation | Enterprise |

`0.5.0`, bootstrap maturity C. Operational endpoints, a Go version
client and a reviewed v1 local operation/receipt contract exist. No
executable-operation implementation, SBRP implementation, remote stop guarantee
or verified robot support is claimed. URLs are target addresses; no package or
binary is published by this initialization.

## Quickstart

Start the core locally using its README (default `127.0.0.1:2046`, PostgreSQL
loopback `2047`). From `sdk/go/`, with Go 1.25+:

```powershell
go test ./...
go vet ./...
go run ./examples/version http://127.0.0.1:2046
```

This source-tree example needs a running local core, but no mutation or keys.
It is not a published-package consumer check. Operational API:
[contracts/bootstrap.openapi.json](contracts/bootstrap.openapi.json). Local
operation/receipt contract:
[contracts/operations/v1/operations.openapi.json](contracts/operations/v1/operations.openapi.json).
Public usage: [docs/README.md](docs/README.md),
[operation contract usage](docs/operations.md) and the
[public Go SDK](docs/sdk.md).

Apache-2.0 covers this repository's original source, not imported core binaries.
Future release bundles must retain AGPL notices, matching corresponding source,
SBOM, checksums and third-party licenses. No binaries are currently included.
