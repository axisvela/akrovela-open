# Changelog

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
