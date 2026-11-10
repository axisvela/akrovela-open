# Changelog

## Unreleased

- Extend the reviewed satellite v1 contract for F4a multi-execution budgets:
  optional `pre_release` execution model, budget mode, ack/exhaustion/capability
  profiles; optional consumption `execution_index`, `repeat_of_attempt_id` and
  `resource_usage`; optional `consumption_ack` per-item results and reason.
  Additive only, fixture-mode evidence scope, no canonical-table change.
- Review post-attempt observation binding and tighten the receipt contract:
  `effect_status=verified` requires a verified, independent verifier, at least
  one accepted observation and `freshness_decision=accepted` on every recorded
  observation. Add `receipt-verified-with-unknown-verifier.json` and
  `receipt-verified-with-stale-observation.json` negative fixtures. No wire
  schema widening.

## 0.1.0 - 2026-11-09

- Bootstrap operational OpenAPI, dependency-free Go client and source example.
- Add the public Go SDK packages `operationsv1` and `satellitev1` with strict JSON
  codecs, plus the experimental `offlineadmissionv1` transport types; schema/Go
  field parity is tested both ways. No network client, signature verification,
  authorization or device trust is included.
- Add reviewed v1 local operation/receipt contract, valid and invalid fixtures,
  and public usage; no executor, database or live KINETOVELA integration.
- Correct the unreleased v1 binding digest to cover resource limits, observation
  references and idempotency retention, so changed effective fields conflict.
- Record that AKROVELA embeds PRAXO-CORE and consumes this contract through
  exact-operation tools without changing the v1 wire schema or SDK.
- Add the reviewed satellite command protocol v1, fixtures and review record;
  fixture-mode only, no live/ROS/effect claim.
- Add the F4a offline pre-release and consumption messages to the satellite
  contract; coordinator-allocated budgets only, no satellite-local admission.
- No published SDK package, core binary or execution implementation is claimed.
