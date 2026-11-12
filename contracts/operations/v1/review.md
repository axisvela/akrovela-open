# Contract Review — Local Operation v1

**Status:** reviewed for integration, 2026-11-08. This is an in-repository
negative-case review, not external customer sign-off, runtime conformance or
security certification.

## Identity

| Item | Value |
| --- | --- |
| Contract | `contracts/operations/v1/operations.openapi.json` |
| Contract line | `v1` |
| Product version | `0.5.0` |
| Consumer | KINETOVELA Site Gateway delegation in fixture mode |
| Maintainer | AKROVELA maintainers |
| Review method | Fresh-context review of contract, fixtures and authority boundaries |

The original review was performed at product `0.1.0`. The additive F5 observation/
effect amendments were reviewed for `0.2.0` without a breaking wire change. No
operation-contract change was made for `0.3.0`; the reviewed contract was carried
forward. No operation-contract change was made for `0.5.0`.

## Ownership

AKROVELA owns local admission, attempts, execution facts, cancellation facts and
receipts. KINETOVELA retains Mission, assignment, lease, fencing and physical
effect authority. PEIRAVELA owns experiment and raw simulation evidence. The
verifier owns the scoped arrival result. Independent safety owns stop and
escalation. No component may substitute its fact for another's authority.

## Negative and Fault Cases

| Case | Required outcome | Coverage |
| --- | --- | --- |
| Missing or expired authentication | Reject with `unauthenticated` | Contract response and review |
| Body tenant/workspace/site mismatch | Reject with `forbidden_scope`; no cross-tenant read | Contract scope binding |
| Unknown operation type or shell/URI field | Reject with `operation_type_not_allowed` or `invalid_request` | Operation allowlist and closed payload |
| Same key, different payload digest | Reject with `operation_key_payload_conflict` | Contract fields; stateful test required |
| Same key, different target/epoch/fence/authority binding | Reject with `operation_key_binding_conflict` | `binding_digest`; stateful test required |
| Same key, different resource limits or idempotency retention | Reject with `operation_key_binding_conflict` | `binding_digest`; stateful test required |
| Same key with non-empty observation refs | Reject with `missing_dependency` until observations are supported | Current deterministic executor rejects non-empty refs before admission; the digest still binds the field |
| Missing operation key or fence | Reject as invalid | `submit-missing-fencing-token.json` |
| Zero fencing token | Reject as invalid | `submit-zero-fence.json` |
| Missing binding digest | Reject as invalid | `submit-missing-binding-digest.json` |
| Operation key contains `/` | Reject as invalid | `submit-slash-key.json` |
| Forged, replayed or expired authority | Reject; no effect | Authority fields and review |
| Replayed single-use approval | Reject; no second effect | Review; implementation gate |
| Stale node boot epoch or fencing token | Reject at admission and effect boundary | `stale_node_boot_epoch`, `stale_fence` |
| Simulator clock used for validity | Reject; validity requires `authority-utc` | `validity-sim-clock.json` |
| Validity expired after queue or approval delay | Reject at final release | Review; implementation gate |
| Resource limit exceeded | Reject with `resource_limit_exceeded` | Closed `ResourceLimits` |
| Missing, stale, replayed, wrong-robot or wrong-operation observation | Block dependent action; `unknown` stays unknown | Observation subject and freshness fields |
| Accepted observation without subject | Reject as invalid | `receipt-observation-missing-subject.json` |
| Cancellation accepted reported as stop | Reject; `stop_confirmed` remains separate | `cancel-stop-without-source.json` |
| False stop with a confirmation source | Reject as invalid | `cancel-false-with-source.json` |
| ROS action termination used as stop proof | Reject; action termination is an execution fact | `action_terminated` and review |
| Verified effect without a receipt | Reject as invalid | `operation-view-verified-without-receipt.json` |
| Receipt missing verifier | Reject as invalid | `receipt-missing-verifier.json` |
| Collapsed operation state | Reject as invalid | `operation-view-missing-cancellation.json` |
| Secret in receipt or evidence | Redact; evidence failure is separate from effect | Review; implementation gate |
| Live verified effect based on simulated evidence | Reject as invalid | `receipt-live-simulated.json` |

## Findings Resolved Before Review

- Cancellation separates request, rejection, unknown and `stop_confirmed`;
  ROS action termination is not stop confirmation.
- `confirmed_no_effect` is a distinct execution outcome; `unknown` is not retry
  permission.
- Idempotency uses payload and binding digests, so a target or epoch swap cannot
  silently return the old operation.
- Validity uses an authoritative clock; simulator time is only an observation
  clock.
- Observations bind source, trust domain, robot, operation, attempt, execution
  epoch, map, target digest, sequence, source epoch, freshness and skew. An
  accepted observation must declare subject, max age, clock skew and source epoch.
- Operation and receipt identifiers are URL-safe and contain no `/`.
- The safety envelope is bound as a reference, not presented as proof of safety.
- Authority and approval revocation and single-use consumption are specified as
  error codes and effect-boundary gates; no behavioral evidence is claimed.
- Receipt types and conditional requirements separate admission receipts from
  execution evidence; a verified effect needs a verified verifier result and an
  observation, and a stop needs a confirmation source and reference.
- Digests use RFC 8785 JCS; the binding digest includes target, epoch, authority,
  resource limits, observation references and idempotency retention, and the server
  recomputes it.
- Zero fencing tokens are rejected by the contract.
- Simulation, ROS action state, receipts and Mission completion remain distinct.

## F5 amendment (post-attempt observations)

This amendment does not widen the v1 wire schema. It makes the existing
receipt conditions explicit and adds one negative fixture.

- Submit `observation_refs` stays empty. Observations are not referenced at
  admission; they bind to the attempt after claim, through `attempt_id` and
  `execution_epoch` in the receipt's observation subject.
- `effect_status=verified` requires `execution_status=completed`,
  `verifier.result=verified`, an independent verifier result and at least one
  accepted observation. Every observation recorded on a verified effect must
  have `freshness_decision=accepted`. Invalid examples:
  `receipt-verified-with-unknown-verifier.json` and
  `receipt-verified-with-stale-observation.json`.
- Missing, stale, replayed, wrong-robot or wrong-operation observations leave
  the effect `unknown`. Failures map to `unknown`, not success, unless a valid
  observation shows the robot outside tolerance. F5a enforces an accepted
  observation with a non-empty `subject`; full observation-envelope, freshness
  arithmetic, clock correspondence and target-digest binding validation remain
  F5b gates.

## Live Receipt Label

`contract_mode=live` is an accepted label. It does not by itself widen the
contract or add runtime evidence. A live `effect_status=verified` receipt still
requires a verified independent verifier and accepted non-simulated observations;
the label does not substitute for that evidence. A live admission receipt
(`receipt_type=admission`, `effect_status=not_evaluated`) records admission only
and does not claim an execution, a verifier result or a physical effect, as shown
by `fixtures/valid/receipt-live-admission.json`. Live KINETOVELA delegation
remains a Phase 3 behavior and is blocked until its owning gates produce evidence.

## Residual Risks

| Risk | Effect | Next gate |
| --- | --- | --- |
| No external consumer sign-off | Contract may not match consumer needs | Consumer review before live integration |
| KINETOVELA Mission submit/query/cancel is not public | Live Mission integration is blocked | KINETOVELA public contract |
| No executor or database implementation | Contract is not runtime evidence | Stage 2 / P1 kernel |
| No behavioral security tests | Fixture review is not security evidence | Deterministic and simulation gates |
| Verifier shares simulation host | Independence is not trust-domain independence | Fault injection and separate evidence |
| Local test validator is a JSON Schema subset | Some semantic coverage is not enforced by the suite | External `jsonschema` and `openapi-spec-validator` checks were run on this revision; add them to CI before relying on the suite alone |
| No stateful idempotency or observation fixtures | Conflict and replay behavior is specified but not executed | Deterministic executor gate |

## Decision

Publish `v1` as a reviewed integration contract. Keep live KINETOVELA integration,
execution implementation, security validation evidence and hardware claims
blocked until their owning gates produce evidence.
