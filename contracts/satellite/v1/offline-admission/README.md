# Offline Admission Transport Candidate

**Status: proposed, awaiting consumer/security review. Not enabled or a released protocol.**
This optional companion to `akrovela.satellite/v1` does not change the existing
`satellite.schema.json` message set. Schema-valid messages are not authorization,
execution evidence or proof of an accepted device profile.

## Upstream Pin

AEGIVELA contract `aegivela.io/offline-admission/v1.0` at
[`684886a1f62e5e8f2159250a584ccd8e4187f667`](https://github.com/axisvela/aegivela-open/tree/684886a1f62e5e8f2159250a584ccd8e4187f667/contracts/offline-admission/v1.0).
`upstream.json` records the 20 consumer case IDs, not test results. All cases
remain not executed. The upstream pin contains stale draft/unpublished prose;
file availability does not authorize deployment or offline enablement.

## Negotiation and Compatibility

1. A new peer requests `offline_admission_v1` in the existing `hello` capabilities.
2. A supporting peer replies with the candidate `hello_ack` containing that exact
   capability only when it was requested and local F4b support is configured.
3. Both peers bind this negotiation to the authenticated session, destination,
   scope and epochs. A legacy ACK without the capability leaves F4b disabled.
4. Only after successful negotiation may the four extension message types be used.
   Do not send the extended ACK to an old peer that did not request it.

Validate the companion schema only for this negotiated extension; continue to use
the original schema for legacy v1 traffic. Capability support is not permission to
execute and cannot substitute for a device-trust provider or signed authority.
The current binaries implement neither this negotiation nor these messages.

## Messages and Authority

| Message | Direction | Meaning |
| --- | --- | --- |
| `offline_grant` | Coordinator to satellite | Relay the exact offline and parent compact artifact bytes with the pinned revision and artifact digest |
| `offline_grant_ack` | Satellite to coordinator | Persisted material or rejected transport; never a grant, admission or effect decision |
| `local_admission_batch` | Satellite to coordinator | Append-only local state assertions, bounded to 100 facts |
| `local_admission_ack` | Coordinator to satellite | Highest contiguous fact sequence durably committed for the indicated original admission epoch |

The coordinator must not reserialize, change audience, re-sign, attenuate or
partition the relayed artifact. It must check destination bindings against
verified claims before relay. The satellite independently verifies the offline
artifact, parent, trust bundle, local submitter, reservation evidence, time and
device profile before any admission. Transport-persisted bytes alone authorize
nothing. Parent keys and offline-signing keys belong to distinct trust domains.
No caller-provided key URL or unsigned scope/budget may enter the relay body.

The 1 MiB encoded frame and 256 queued-command bounds remain unchanged. The
schema's 100-fact bound is not a byte-size guarantee; the transport must enforce
the frame limit before decoding and split batches accordingly. Neither artifact
strings nor unrestricted arguments belong in telemetry or evidence logs.

## Scope and Epoch Mapping

- Derive tenant from verified AEGIVELA namespace through configured mapping.
  Envelope tenant/workspace/site are comparison-only, not identity selectors.
- Match `device.deviceRef` to `satellite_id`/`node_id`, `device.sourceEpoch` to
  the trusted boot epoch, and signed workspace/site to configured scope.
- `source_epoch` in the envelope is the authenticated transport source's current
  epoch. `admission_source_epoch` in a fact batch/ACK is the original epoch of
  the recorded work. All facts in a batch belong to that original epoch.
- Historical epochs may be reconciled only against registered device history and
  the original grant; they never authorize new work after reboot. Larger integers
  alone cannot create trusted rebinding. An ACK for one epoch cannot clear another.

## Digests and Request Identity

| Field | Bytes hashed / domain |
| --- | --- |
| `artifact_digest` | Lowercase bare SHA-256 hex of the exact compact offline artifact ASCII bytes |
| `authority_payload_digest` | AEGIVELA SHA-256 hex over JCS of `{operationType, payload}` |
| `authority_binding_digest` | AEGIVELA bare SHA-256 hex over the exact binding projection defined by its pinned canonicalization contract |
| `operation_payload_digest` | Existing AKROVELA `sha256:` operation-payload digest; not the AEGIVELA payload projection |
| `operation_binding_digest` | Existing AKROVELA `sha256:` digest over all immutable effective operation fields |

Do not strip a prefix and assume that two digest domains are equivalent. Verify
each independently. The AEGIVELA binding projection includes namespace, audience,
action, resource, structured-resource digest, task, scope, device/key/epoch,
parent reference/digest and validity. Scope array order follows the signed bytes;
no implicit sorting or set normalization before hashing.

The unique logical admission key is namespace/device/original-epoch/operation-key/
attempt-ID. Persist and compare immutable request/binding/payload digests. Same
key with identical content returns its durable state without spending or starting
again; changed content rejects. Putting changing digests into the uniqueness key
must not permit another admission or spend.

## Units and Reservation Facts

AEGIVELA uses integer milliseconds, millimetres and millimetres/second. Convert
AKROVELA metre-valued bounds from their canonical decimal representation using
exact decimal arithmetic and a conservative ceiling after multiplying by 1000;
never round down, truncate a float or accept NaN/infinity/negative/overflow.
Reject a conversion above 9007199254740991. For example 0.0011 m reserves 2 mm,
not 1 mm; 0.45 m/s requires a 450 mm/s ceiling.

Reserve the full worst-case count/duration/goal/path bounds before execution.
Speed and payload size are per-operation ceilings, not additive pools. Each
admission reserves one operation; no automatic retries or synthetic attempts.
`reserved_resources` on each fact is the immutable reservation snapshot, not a
new charge on each state transition. Repeated facts must not spend again.

`fact_id` and its sequence identify an immutable transition assertion. New states
append new facts under the same admission; never overwrite an earlier unknown
assertion. `terminal` is a local executor lifecycle assertion, not physical-effect
or canonical-stop proof. Existing receipts retain their separate meanings.
Mission/assignment/lease/fencing remain required trusted product constraints;
none may be fabricated from this artifact.

## ACK and Recovery Rules

Fact sequencing/ACK is scoped by configured tenant/workspace/site, namespace,
satellite and original admission epoch, not by a transient TCP session. Validate
the batch namespace against every fact, verified grant and configured mapping.
The ACK repeats that namespace; it cannot acknowledge another namespace's stream.
Validate all facts and their immutable identity/digest before atomically committing. ACK
only the contiguous committed prefix. On a rejected row or sequence gap, do not
advance past it; duplicates must match the stored fact exactly.

Unknown or quarantined work keeps its reservation. Expiry, cancellation,
revocation, an ACK or reconnect is not evidence of non-use and does not refund
authority. Authenticate and deduplicate late facts without renewing authority.
Drain pending facts before new online or offline work after reconnect.

F4a/online and F4b must refer to an authority-linked non-overlapping partition;
until that relation is proven, mixed use of the same allowance must reject.
Two independently valid local counters cannot create two spending pools.

## Unresolved Security Review Questions

- `no-reservation`: what authenticated evidence lets a disconnected consumer
  prove that its allowance was committed by AEGIVELA? The pinned artifact has no
  explicit reservation reference/proof field. Do not invent a signed proof or
  treat a scope-bearing signature as a committed reservation.
- `parent-revoked-disconnected`: what observation or elapsed deadline supports
  the expected rejection? The documented finite disconnected policy does not
  deliver instantaneous knowledge of a remote revocation. Reject known revoked,
  connected-only or expired-freshness authority, but do not claim knowledge of an
  undelivered event. This case requires owner clarification before acceptance.

Owner review of this candidate and these two questions is required before a
dependent runtime implementation. Device rollback/clone and single-consumer
proofs also remain separate acceptance prerequisites; test doubles cannot supply
production proof.

## Fixture Limits

The five valid fixtures are transport/schema examples only. Their short compact
tokens deliberately have no valid authorization signature. The batch includes
example digest values from different projections, not an internally consistent
signed request. They must fail runtime authorization. Runtime tests must use
verified upstream vectors and authenticated test facts instead.

Run `python -m unittest discover -s tests -v` from the public repository root.
The dependency-free test validator checks the schema features used here, not the
entire JSON Schema vocabulary or cross-message semantic invariants.
