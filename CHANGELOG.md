# Changelog

## 0.11.0 - 2026-11-16

- Align the public version with core `0.11.0`; the reviewed wire contracts are
  carried forward.
- Core `0.11.0` adds the F5 Linux simulation prerequisites: a pinned stack
  manifest and provisioner, a JSONL Gazebo observation exporter with
  attempt-start binding, a Linux preflight gate and an opt-in cross-product CI
  E2E. These are component increments, not real ROS 2/Gazebo/Linux or
  PEIRAVELA acceptance.
- No public operation, satellite or delegation wire change is added.

## 0.10.0 - 2026-11-15

- Align the public version with core `0.10.0`; the reviewed wire contracts are
  carried forward.
- Core `0.10.0` adds the TPM trust-profile approval gate, a local device
  inventory whitelist, approval expiry/revocation enforcement, policy content
  validation in both preflight commands, simulated counter-floor rejection and
  expanded simulated F5 fault paths.
- No public operation, satellite or delegation wire change is added.

## 0.9.0 - 2026-11-15

- Add `operationsv1.Client.Wait` to poll until execution reaches a terminal
  status. It does not wait for effect verification or retry effects.
- Add `operationsv1.Client.SubmitAndWait` to submit and poll to a terminal
  execution status.
- Reject empty operation keys, digests, cancel request IDs and receipt IDs before
  sending a request.
- Add `-wait` to the runnable operations example submit command.
- Align the public version with core `0.9.0`; the reviewed wire contracts are
  carried forward.

## 0.8.0 - 2026-11-15

- Add an `operationsv1` HTTP client for operation submit/query/cancel/receipts.
  It maps the reviewed error body and requires caller-supplied digests; it does
  not retry effects or confirm a stop.
- Add `PayloadDigest`, `BindingDigest` and `FillDigests` JCS helpers to the
  `operationsv1` SDK. They use RFC 8785 numeric and string canonicalization and
  are checked against the reviewed operation-binding fixture.
- Add a runnable `examples/operations` command for submit, query, cancel and
  receipt calls through the public `operationsv1` client.
- Align the public version with core `0.8.0`; the reviewed operation, satellite
  and delegation wire contracts are carried forward.

## 0.7.0 - 2026-11-15

- Align the public version with core `0.7.0`. The reviewed v1 operation,
  satellite and delegation contracts are carried forward; there is no public
  operation, receipt, satellite or delegation wire change.
- Core `0.7.0` adds TPM-backed offline device evidence, one-time AK/EK
  enrollment binding, a TPM-resident device signing key, one-time verifier
  challenges, Roughtime, offline deployment preflight and satellite
  offline-runtime grant/parent-fact persistence.
- This alignment adds no production device trust, offline guarantee, ROS 2,
  hardware or physical-effect claim.

## 0.6.0 - 2026-11-13

- Align the public version with core `0.6.0`. The reviewed v1 operation,
  satellite and delegation contracts are carried forward; there is no public
  operation, receipt, satellite or delegation wire change.
- Core `0.6.0` adds PostgreSQL Mission control-plane persistence, cursor
  pagination and durable delegation state; the KINETOVELA cursor contract is
  published in `KINETOVELA-open`.
- No production identity, offline guarantee, ROS 2, hardware or physical-effect
  claim is added.

## 0.5.0 - 2026-11-11

- Align the public version with core `0.5.0`. The reviewed v1 operation,
  satellite and delegation contracts are carried forward; core adds the
  KINETOVELA live delegation components without a public wire change to the
  existing operation/receipt contract.
- Add the reviewed `v1` delegation contract for the KINETOVELA-to-AKROVELA
  bounded operation envelope: schema, valid `fixture`/`live` fixtures, negative
  fixtures and public usage. The current adapter accepts `fixture` mode only;
  `live` is published and reviewed for Phase 3 but not accepted yet.
- Add the live admission receipt fixture to the operation contract and record the
  live receipt label in the operation review.
- Add the delegation review record, contract index, public docs and changelog
  entry. No live integration, real verifier or physical effect is claimed.

## 0.4.0 - 2026-11-11

- Align the public version with core `0.4.0`. No public contract or SDK behavior
  change; the reviewed v1 operation and satellite contracts are carried forward.
  Core `0.4.0` adds F5b verifier and bridge state durability, optional verifier
  shared-token auth and service-level fake end-to-end coverage without a public
  wire change.

## 0.3.0 - 2026-11-11

- Align the public version with core `0.3.0`. No public contract or SDK behavior
  change; the reviewed v1 operation and satellite contracts are carried forward.
  Core `0.3.0` adds the F5b arrival verifier process and cancel/restart
  reconciliation without a public wire change.

## 0.2.0 - 2026-11-10

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

