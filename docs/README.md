# Public Usage

`GET /healthz` reports process liveness; `GET /readyz` reports bootstrap database
and schema availability; `GET /version` reports product/version/bootstrap stage.
These are unauthenticated local-development endpoints, not execution readiness
or governance certification. The default API port is 2046.

The SDK only reads version. Operation execution and enterprise capabilities remain
unimplemented. The v1 local operation/receipt contract is published in
`contracts/operations/v1/operations.openapi.json`; see
[operation contract usage](operations.md). It is a reviewed contract, not an
implemented executor. The root `VERSION` and OpenAPI info version align at
`0.2.0`.
Go modules start without a major suffix while 0.x; a future API major follows Go
module versioning. Publishing the nested SDK needs `sdk/go/vX.Y.Z`, separately
from product root tags; no such tags or registry availability are claimed here.
