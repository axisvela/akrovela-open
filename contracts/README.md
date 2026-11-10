# Public Contracts

| Contract | Status | Scope |
| --- | --- | --- |
| [bootstrap.openapi.json](bootstrap.openapi.json) | `0.3.0` | Health, readiness and version only |
| [operations/v1/operations.openapi.json](operations/v1/operations.openapi.json) | `v1`, reviewed | Local operation submit/query/cancel/receipt protocol |

The operation contract is not an implementation. Review the
[contract review](operations/v1/review.md) and
[valid/negative fixtures](operations/v1/fixtures/) before integration. Breaking
changes require a new contract line and compatibility review.
