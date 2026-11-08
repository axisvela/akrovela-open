# Contract Review - Satellite Command Protocol v1

**Status:** reviewed for integration, 2026-11-08. This is an in-repository review,
not external sign-off, runtime conformance or a security certification.

## Identity

| Item | Value |
| --- | --- |
| Contract | `contracts/satellite/v1/satellite.schema.json` |
| Contract line | `v1` |
| Protocol version | `akrovela.satellite/v1` |
| Product version | `0.1.0` |
| Mode | Deterministic fixture only |

## Approved Decisions

1. The satellite owns a local execution-assertion ledger. The coordinator mints
   canonical receipts on upload, keyed by `(satellite_id, satellite_receipt_id)`.
2. A satellite reports a local stop assertion. Canonical `stop_confirmed` requires
   an independent observation or safety authority.
3. `source_epoch` equals the durable device boot epoch (`node_boot_epoch`).
   Process restart does not advance it.
4. Restart after durable acceptance executes only when the binding, epochs,
   dispatch and validity windows remain valid and no effect was issued; otherwise
   the outcome stays `unknown`.
5. Only `lease_token` dispatch fencing is in scope.
6. This contract is fixture-mode only; it makes no live operation/effect claim.
7. First binding is length-prefixed JSON: max frame 1 MiB, queue 256, receipt
   batch 100.
8. First implementation is Go.

## Ownership

AKROVELA owns coordinator admission, epoch, command journal and canonical
receipts. KINETOVELA retains Mission, assignment, lease, fencing and effect
authority. The satellite reports scoped assertions. Independent verification and
safety remain separate owners.

## Negative and Fault Cases

| Case | Required outcome |
| --- | --- |
| Unknown protocol version or message type | Reject; no fallback |
| Unknown field | Reject (`additionalProperties: false`) |
| Missing scope, satellite or coordinator identity | Reject |
| Missing source/coordinator epoch | Reject |
| Missing fence or dispatch expiry | Reject |
| Expired dispatch window | Reject before any effect |
| Lower coordinator or source epoch | Reject |
| Same operation/attempt with a different binding digest | Conflict; do not replace |
| Duplicate message/operation | Return persisted result; do not re-execute |
| Cancel accepted reported as stop | Reject; `local_stop_asserted` is not canonical stop |
| Same satellite receipt ID with a different digest | Equivocation conflict |
| Receipt replay | Deduplicate by satellite receipt ID |
| Process restart with unchanged source epoch | Resume, do not re-execute an effect |
| Device reboot with changed source epoch | Require explicit rebinding |
| Wrong identity/audience | Reject |
| Oversized frame or full queue | Reject; no partial execution |

## Residual Risks

| Risk | Effect | Next gate |
| --- | --- | --- |
| No external consumer sign-off | Contract may not match the first satellite | Consumer review before live use |
| No implementation | Contract is not runtime evidence | Phase 1-3 implementation |
| Frame/queue limits are product choices | May need tuning | Measured Linux/Windows satellite |
| Fixture-only | No live operation or effect evidence | ROS/simulation gate |
| Key/identity provider is open | Cannot claim production identity | EIDOVELA/session-key decision |

## Decision

Publish `v1` as a reviewed fixture-mode integration contract. Live operation,
ROS/Nav2, offline admission, MCU bindings, effect verification and hardware claims
remain blocked until their owning gates produce evidence.
