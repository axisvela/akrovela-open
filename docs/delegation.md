# KINETOVELA Delegated Operation Contract v1

**Status:** public contract, product `0.4.0`, contract `v1`. The contract is
reviewed for integration; it is not an implementation. No live KINETOVELA
integration, ROS 2 execution or physical effect is claimed.

Canonical contract:
[delegated-operation.schema.json](../contracts/delegation/v1/delegated-operation.schema.json).

Fixtures:
[valid and invalid cases](../contracts/delegation/v1/fixtures/).

Review record:
[contract review](../contracts/delegation/v1/review.md).

## Purpose

The delegation contract is the bounded envelope KINETOVELA uses to hand one
local operation to AKROVELA. It carries the Mission, assignment, lease, fencing,
authority, safety and approval references for one `ros2.nav2.navigate_to_pose@1`
operation, together with its target, payload, validity, resource limits and
idempotency retention. AKROVELA maps the envelope onto the
[local operation contract](operations.md) and remains authoritative for local
admission, execution and receipts.

## Modes

| `contract_mode` | Meaning | Current adapter |
| --- | --- | --- |
| `fixture` | Deterministic, non-live assembly of the local operation request | Accepted |
| `live` | Published and reviewed Phase 3 label | Rejected with `missing_dependency` |

`fixture` is the only mode the current core delegation adapter accepts. `live` is
published and reviewed for Phase 3 but is not accepted until the Phase 3 server
gate exists. Accepting the label would not itself add runtime evidence.

## Authority

KINETOVELA owns Mission, assignment, lease, fencing and physical effect.
AKROVELA owns local admission, attempts, execution facts, cancellation facts and
receipts. A delegated operation does not transfer Mission authority, and neither a
delegated operation nor a receipt completes a Mission. PEIRAVELA owns experiment
and raw simulation evidence; the verifier owns the scoped arrival result; and
independent safety owns stop and escalation.

## Required Fields

A delegation requires:

- `schema_version`, fixed to `akrovela.delegated-operation/v1`;
- `contract_mode`, `delegation_id` and `operation_key`;
- `operation_type`, fixed to `ros2.nav2.navigate_to_pose@1`;
- `scope` (`tenant_id`, `workspace_id`, `site_id`);
- `mission_ref`, `assignment_ref`, `lease_ref` and `fencing_token`;
- `authority_ref`, `safety_envelope_ref` and `approval_ref`;
- `expected_effect`, `target` and `payload`;
- `observation_refs` (empty for the current adapter);
- `validity`, `resource_limits` and `idempotency_retention_seconds`.

The delegation envelopes carry references; they do not embed grants, approvals or
safety proof.

## Operation Key and Digests

`operation_key` is the stable, URL-safe identity of the local operation. The
delegation envelope does not carry digests: the AKROVELA adapter computes them
when it builds the local submission.

- `payload_digest` is SHA-256 over RFC 8785 JSON Canonicalization Scheme (JCS) of
  `{operation_type, payload}`.
- `binding_digest` is SHA-256 over JCS of the effective binding: `scope`,
  `target`, `mission_ref`, `assignment_ref`, `lease_ref`, `fencing_token`,
  `safety_envelope_ref`, `authority_ref`, `approval_ref`, `validity`,
  `resource_limits`, `observation_refs` and `idempotency_retention_seconds`.

Reusing an `operation_key` with a different payload rejects as
`operation_key_payload_conflict`; a different target, epoch, fence, authority or
validity binding rejects as `operation_key_binding_conflict`. Neither creates an
effect. The server recomputes both digests and rejects a mismatch.

## Fixtures and Review

The [valid fixtures](../contracts/delegation/v1/fixtures/valid/) show a
`fixture` delegation and the published `live` envelope. The
[invalid fixtures](../contracts/delegation/v1/fixtures/invalid/) cover a missing
delegation id, a wrong schema version, a bad contract-mode value, a zero fence, a
missing mission reference and a missing expected effect.

The [review record](../contracts/delegation/v1/review.md) lists the ownership
boundaries, the mode limits and the negative and fault cases.

## Limits

- No live KINETOVELA integration, Mission submit/query/cancel surface or server
  gate exists here.
- No ROS 2 execution, executor implementation or real verifier is provided.
- No physical effect is claimed.
- The `live` label is published and reviewed for Phase 3, not accepted by the
  current adapter.
- A receipt is local evidence, not Mission completion or independent assurance
  certification.
