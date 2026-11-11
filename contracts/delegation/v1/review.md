# Contract Review - Delegated Operation v1

**Status:** reviewed for integration, 2026-11-11. This is an in-repository review,
not external sign-off, runtime conformance or live integration evidence.

## Identity

| Item | Value |
| --- | --- |
| Contract | `contracts/delegation/v1/delegated-operation.schema.json` |
| Contract line | `v1` |
| Product version | `0.4.0` |
| Consumer | KINETOVELA live delegation |
| Maintainer | AKROVELA maintainers |
| Review method | In-repository review of contract, fixtures and ownership boundaries |
| Reviewed `akrovela` core revision | `af16a010ae299f01247d8a5c389fd17815bc721e` |
| Reviewed `akrovela-open` revision | `8f01dbd9f5589e63310255a07a5d7ed2ee258fbc` |
| Upstream KINETOVELA core revision | `545023997bff5f6287a80fe5817c1782e52ae570` |

The contract line is `v1`; the product version tracks the `0.4.0` release. The
reviewed core revision is the `akrovela` source tree that contains the controlled
delegation adapter. The delegation references the KINETOVELA Mission and effect
contracts reviewed at KINETOVELA core `545023997bff5f6287a80fe5817c1782e52ae570`;
AKROVELA does not import KINETOVELA internals. This record is evidence for an
in-repository review only.

## Ownership

KINETOVELA owns Mission, assignment, lease, fencing and physical effect.
AKROVELA owns local operation admission, attempts, execution facts, cancellation
facts and receipts. The delegation envelope carries KINETOVELA references into
the local operation contract; it does not transfer Mission authority. A delegated
operation or its receipt never completes a Mission. PEIRAVELA owns experiment and
raw simulation evidence, the verifier owns the scoped arrival result, and
independent safety owns stop and escalation. No component may substitute its fact
for another's authority.

## Modes

`fixture` is the only contract mode the current core delegation adapter accepts.
A fixture delegation is a deterministic, non-live assembly of the local operation
request; it claims no live KINETOVELA integration, real verifier or physical
effect.

`live` is published and reviewed for Phase 3, but the current core adapter does
not accept it. The adapter rejects `contract_mode=live` with
`missing_dependency` until the Phase 3 server gate exists. Reviewing the label
does not widen the wire schema and does not provide live integration evidence.

## Negative and Fault Cases

| Case | Required outcome | Coverage |
| --- | --- | --- |
| Missing delegation id | Reject as invalid | `delegated-operation-missing-delegation-id.json` |
| Wrong schema version | Reject as invalid | `delegated-operation-wrong-schema-version.json` |
| Bad contract mode value | Reject as invalid | `delegated-operation-bad-contract-mode.json` |
| Zero or negative fence | Reject as invalid | `delegated-operation-zero-fence.json`; adapter rejects `< 1` |
| Missing mission reference | Reject as invalid | `delegated-operation-missing-mission-ref.json` |
| Missing assignment reference | Reject as invalid | Schema pattern and adapter reference check |
| Missing lease reference | Reject as invalid | Schema pattern and adapter reference check |
| Missing safety-envelope reference | Reject as invalid | Schema pattern and adapter reference check |
| Missing authority reference | Reject as invalid | Schema pattern and adapter reference check |
| Missing approval reference | Reject as invalid | Schema pattern and adapter reference check |
| Missing expected effect | Reject as invalid | `delegated-operation-missing-expected-effect.json` |
| Oversized expected effect | Reject as invalid | Schema `maxLength` of 512; adapter length check |
| Scope mismatch with the authenticated principal | Reject with `forbidden_scope`; no cross-tenant admission | Schema scope and adapter authorization check |
| Non-empty observation references | Reject with `missing_dependency` | Current adapter rejects non-empty refs before admission |
| Receipt treated as Mission completion | Reject; the receipt is local evidence, not Mission completion | Ownership section and operation contract |
| Accepted cancellation treated as stop | Reject; `stop_confirmed` remains separate | Operation contract and ownership section |
| Remote `$ref` in the schema | Reject; only local `#/` references are allowed | `test_local_refs_only` |
| Live mode without the Phase 3 server gate | Reject with `missing_dependency` | Adapter rejects every non-`fixture` mode |

## Residual Risks

| Risk | Effect | Next gate |
| --- | --- | --- |
| No external consumer sign-off | Contract may not match consumer needs | KINETOVELA review before live integration |
| No live KINETOVELA Mission contract | Live delegation is blocked | KINETOVELA public Mission and effect-admission surfaces |
| No runtime conformance or live integration | Fixture and review evidence only | Phase 3 server gate and live integration evidence |
| No real verifier or physical effect | No execution or effect claim | Owning verifier and execution gates |
| Local test validator is a JSON Schema subset | Some semantic coverage is not enforced by the suite | External `jsonschema` checks before relying on the suite alone |

## Decision

Publish `v1` as a reviewed delegation contract for fixture-mode integration. Treat
`live` as a published Phase 3 label that the current adapter rejects. Keep live
KINETOVELA integration, real verifier evidence and physical-effect claims blocked
until their owning gates produce evidence.
