# Public Contracts

| Contract | Status | Scope |
| --- | --- | --- |
| [bootstrap.openapi.json](bootstrap.openapi.json) | `0.4.0` | Health, readiness and version only |
| [operations/v1/operations.openapi.json](operations/v1/operations.openapi.json) | `v1`, reviewed | Local operation submit/query/cancel/receipt protocol |
| [delegation/v1/delegated-operation.schema.json](delegation/v1/delegated-operation.schema.json) | `v1`, reviewed | KINETOVELA-to-AKROVELA delegated operation |

The operation contract is not an implementation. Review the
[contract review](operations/v1/review.md) and
[valid/negative fixtures](operations/v1/fixtures/) before integration. Breaking
changes require a new contract line and compatibility review.

The delegation contract is the bounded KINETOVELA-to-AKROVELA envelope for one
local operation. Review the [delegation review](delegation/v1/review.md) and its
[fixtures](delegation/v1/fixtures/). The current adapter accepts `fixture` mode
only; `live` is published and reviewed for Phase 3 but is not accepted yet.
