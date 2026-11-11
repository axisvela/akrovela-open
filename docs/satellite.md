# Satellite Command Protocol v1

**Status:** public contract, product `0.4.0`, protocol `v1`, fixture mode only.
This is not an implementation and claims no live operation, ROS, effect or
hardware behavior.

Canonical contract:
[satellite.schema.json](../contracts/satellite/v1/satellite.schema.json).

Fixtures:
[valid and invalid cases](../contracts/satellite/v1/fixtures/).

Review record:
[contract review](../contracts/satellite/v1/review.md).

## Scope

The contract carries one coordinator's bounded fixture command to one execution
satellite and the satellite's scoped assertions back. It does not create a
Mission, assignment, effect, world state or safety decision.

## Envelope

Every message carries `protocol_version`, `message_type`, `message_id`,
`session_id`, `sequence`, tenant/workspace/site, `satellite_id`,
`source_epoch`, `coordinator_id`, `coordinator_epoch`, `sent_at` and `body`.

`satellite_id` is the AKROVELA `node_id`. `source_epoch` is the durable device boot
epoch and equals `node_boot_epoch`. Process restart does not advance it.

## Messages

| Message | Direction | Meaning |
| --- | --- | --- |
| `hello` / `hello_ack` | both | Version, identity and capability negotiation |
| `resume` / `resume_result` | both | Signed resume and drain-before-new-work |
| `submit` / `submit_result` | coordinator / satellite | Durable acceptance is separate from transport ACK |
| `status` / `status_result` | both | Authoritative command state |
| `cancel` / `cancel_result` | both | Per-cancel-ID request and local assertion |
| `receipt_batch` / `receipt_ack` | satellite / coordinator | Receipt upload and cumulative acknowledgement |
| `heartbeat` | both | Liveness only; not execution or stop evidence |

Unknown versions, message types or fields reject. There is no silent fallback.

## Effecting Commands

`submit` carries the operation key, attempt identity, payload and binding digests,
target, Mission/assignment/lease references, fence, safety/authority/approval
references, validity, resource limits, empty observation refs, retention,
`dispatch_not_after` and `contract_mode=fixture`.

`dispatch_not_after` is absolute UTC and must not exceed the coordinator leadership
lease. The earlier of dispatch expiry and operation validity wins. An expired
command never starts an effect.

## Idempotency and Conflicts

The dedupe key is `(satellite_id, operation_key, attempt_id, binding_digest)`. A
duplicate returns the persisted result, including terminal results. A different
binding digest, fence or epoch for the same operation/attempt conflicts and does
not replace the earlier command.

## Cancellation

Each `cancel_request_id` has its own status. A satellite may report
`cancel_requested`, `cancel_rejected`, `local_stop_asserted` or `cancel_unknown`.
Canonical `stop_confirmed` requires an independent observation or safety
authority; action termination, process exit, disconnect or heartbeat loss is not a
stop.

## Receipts

The satellite persists append-only local receipts. The coordinator imports them
transactionally into canonical receipts keyed by
`(satellite_id, satellite_receipt_id)`, deduplicates, acknowledges by sequence and
maps the canonical per-operation sequence. A repeated satellite receipt ID with a
different digest is an equivocation conflict.

## Transport

Bindings are reviewed separately: Unix domain sockets, Windows named pipes and
mutually authenticated TLS. Every binding provides framing, size limits,
canonical signed bytes, identity, replay protection and downgrade rejection.
Listening on a socket is not authentication.

## Limits

Maximum frame 1 MiB, command queue 256, receipt batch 100. These are product
settings, not safety limits.

## Offline Admission Candidate

The [F4b companion mapping](../contracts/satellite/v1/offline-admission/README.md)
and [candidate schema](../contracts/satellite/v1/offline-admission.schema.json)
are proposed, not reviewed or enabled. They pin AEGIVELA offline-admission v1.0
and require explicit capability negotiation; they do not widen the existing v1
message set. The 20 upstream consumer acceptance cases have not been executed.
