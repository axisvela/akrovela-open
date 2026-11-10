# Local Operation and Receipt Contract v1

**Status:** public contract, product `0.2.0`, contract `v1`. The contract is
reviewed for integration; it is not an implementation and no executor, database
schema or live KINETOVELA integration is claimed.

Canonical contract:
[operations.openapi.json](../contracts/operations/v1/operations.openapi.json).

Fixtures:
[valid and invalid cases](../contracts/operations/v1/fixtures/).

Review record:
[contract review](../contracts/operations/v1/review.md).

## Scope

The contract covers one bounded local operation on one AKROVELA node:

| Operation type | Meaning |
| --- | --- |
| `ros2.nav2.navigate_to_pose@1` | Move one robot to a pose on a named map through ROS 2/Nav2 |

AKROVELA is authoritative for local operation facts: admission, attempts,
execution, cancellation and receipts. KINETOVELA retains Mission, assignment,
lease, fencing and physical-effect authority. A receipt does not complete a
Mission. ROS topics, Nav2 action results and Gazebo facts are not authority.

## Surfaces

| Surface | Method | Meaning |
| --- | --- | --- |
| `/v1/operations` | `POST` | Submit one operation; idempotent on the key and both digests |
| `/v1/operations/{operation_key}` | `GET` | Query persisted authoritative local state |
| `/v1/operations/{operation_key}/cancel` | `POST` | Record a cancellation request; not a stop |
| `/v1/operations/{operation_key}/receipts` | `GET` | List append-only receipts |
| `/v1/receipts/{receipt_id}` | `GET` | Get one immutable receipt |

Keys are URL-safe and contain no `/`. Authentication uses a bearer principal;
body scope is checked against the authenticated tenant, workspace and site.

## Embedding and Observation Limits

AKROVELA embeds PRAXO-CORE as a product component. Its `praxocore` tools
(submit/status/cancel/receipts) consume this public contract for exactly one
trusted operation assembly; model tool arguments carry no operation, target,
authority or resource selection. Direct non-Agent operation paths remain
supported. Current deterministic runtime execution supports only empty observation
references and no independently verified trajectory. The schemas alone do not
implement grants, approval consumption, resource enforcement or observation
authenticity.

Before implementing prerequisite observations, review their binding separately
from post-attempt evidence: an attempt identity cannot be invented before the
execution authority assigns it. Review exact digest/approval coverage of resource
limits, observation policy and retention rather than assuming every field is in
the current binding digest. This note changes no schema or digest definition.

Keep validity UTC, observation clock/error bounds and simulator epoch/ticks
distinct. Given combined clock-error bound U, require worst-case observation age
at release to satisfy `release_time - observed_at + U <= max_age`; unknown bounds
block dependent actions. Source epoch changes need authenticated rebinding;
repeated reads do not refresh time. A local verified result remains scoped
evidence, not KINETOVELA Mission completion or independent assurance certification.

## Bindings

A submission binds:

- `operation_key` and `protocol_version`;
- `payload_digest` as SHA-256 over RFC 8785 JSON Canonicalization Scheme (JCS) of
  `{operation_type, payload}`;
- `binding_digest` as SHA-256 over JCS of `{scope, target, mission_ref,
  assignment_ref, lease_ref, fencing_token, safety_envelope_ref, authority_ref,
  approval_ref, validity, resource_limits, observation_refs,
  idempotency_retention_seconds}`;
- `node_boot_epoch` and `fencing_token` (fence values start at 1);
- a validity window on an authoritative clock (`authority-utc`);
- resource limits;
- observation references and their freshness decisions.

The server recomputes both digests and rejects a mismatch as
`payload_digest_mismatch` or `binding_digest_mismatch`. Reusing an operation key
with a different payload rejects as `operation_key_payload_conflict`; a different
target, epoch, fence, authority or validity binding rejects as
`operation_key_binding_conflict`. Neither creates an effect.

## State Separation

Admission, execution, cancellation and verified effect are separate axes:

| Axis | Values |
| --- | --- |
| `admission_status` | `accepted`, `rejected` |
| `execution_status` | `not_started`, `started`, `completed`, `failed`, `confirmed_no_effect`, `unknown` |
| `cancellation_status` | `not_requested`, `requested`, `stop_confirmed`, `rejected`, `unknown` |
| `effect_status` | `not_evaluated`, `verified`, `rejected`, `unknown` |

`execution_status=completed` is not `effect_status=verified`.
`confirmed_no_effect` requires positive proof that no effect was issued; a lost
response is `unknown`. `unknown` is never silently retried or relabeled.

Receipts declare a `receipt_type`: `admission`, `execution`, `cancellation`,
`evidence` or `repair`. An admission receipt does not need executor or verifier
facts. An execution, evidence or repair receipt requires an attempt, executor
facts, verifier and observations. `effect_status=verified` requires a verified
verifier result and at least one observation. `stop_confirmed` requires a
confirmation source and reference.

## Cancellation

A cancellation request has its own identity. An accepted request is not a stop.
`stop_confirmed` is true only when a separate safety authority or independent
observation confirms the stop. ROS 2 action termination is an execution fact
named `action_terminated`, not stop confirmation. A confirmed stop does not erase
earlier effects and does not complete a Mission.

## Observations and Verification

Observation references carry source, trust domain, state label, clock domain,
source epoch, sequence, timestamps, digest, freshness bound and subject binding
(robot, operation, attempt, execution epoch, map and target digest). An accepted
observation must declare its subject, max age, clock-skew bound and source epoch.
Simulated observations remain `state_label=simulated`; a live receipt with
`effect_status=verified` cannot rely on simulated evidence. The verifier, not the
executor, makes the final freshness and arrival decision. A passing simulation is
not hardware evidence.

## Security Validation

The review record lists the negative and fault cases that must reject or block:
authentication and tenant mismatch, digest mismatch, operation allowlist,
payload/binding conflict, forged, replayed or revoked authority, approval
consumption or expiry, stale epoch or fence, validity and clock violations,
resource limits, stale or wrong-robot observations, unconfirmed cancellation,
secret redaction, cross-tenant receipt access and simulation/ROS conflation.

The error enum names `payload_digest_mismatch`, `binding_digest_mismatch`,
`authority_revoked`, `approval_consumed` and `approval_expired`. Revocation and
atomic single-use consumption are specified for the effect-boundary
implementation gate; no behavioral security evidence, live service test or
production claim is made.

## Authority and Ownership

| Concern | Owner |
| --- | --- |
| Mission, assignment, lease, fence and effect lifecycle | KINETOVELA |
| Local operation admission, attempts, cancellation facts and receipts | AKROVELA |
| Experiment, perturbations and raw simulation evidence | PEIRAVELA |
| Independent arrival verification | Consumer verifier |
| Local safety, stop and escalation | Independent safety owner |

The first bounded consumer is a KINETOVELA Site Gateway delegation in fixture
mode. It cannot create or complete a Mission and cannot import private KINETOVELA
schemas. Live integration remains blocked until KINETOVELA publishes its Mission
submit/query/cancel and effect-admission surfaces.

## Not Implemented

No operation executor, PostgreSQL operation schema, ROS 2 adapter, live
KINETOVELA adapter, SDK operation methods or behavioral test suite exists in this
repository. The contract and fixtures are the Stage 1 public surface. A standalone
receipt is a fact record keyed by `operation_key`; query the operation view for
the full submitted binding when needed.

## Versioning

`v1` is the contract line; `info.version` tracks the product version. Breaking
changes require a new contract line and ordinary compatibility review. A passing
fixture check is not runtime conformance evidence.
