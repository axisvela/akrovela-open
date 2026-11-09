import copy
import hashlib
import sys
from pathlib import Path
import unittest

sys.path.insert(0, str(Path(__file__).resolve().parent))
from test_operations_contract import SchemaError, load, validate

ROOT = Path(__file__).resolve().parents[1]
DIRECTORY = ROOT / 'contracts/satellite/v1/offline-admission'
SCHEMA = DIRECTORY.parent / 'offline-admission.schema.json'
REVISION = '684886a1f62e5e8f2159250a584ccd8e4187f667'
CASE_IDS = {
    'valid-admit', 'unknown-key', 'audience-mismatch', 'device-epoch-mismatch',
    'operation-not-allowed', 'clock-deadline', 'no-reservation',
    'concurrent-final-credit', 'identical-retry', 'conflicting-retry',
    'persistence-failure', 'uncertain-commit', 'crash-before-start-intent',
    'crash-after-start-intent', 'restart-started-not-terminal',
    'drain-before-new-work', 'journal-rollback', 'device-clone',
    'parent-revoked-disconnected', 'late-receipt-after-revocation',
}


class OfflineAdmissionContractTests(unittest.TestCase):
    def setUp(self):
        self.assertTrue(SCHEMA.is_file(), 'F4b candidate schema is missing')
        self.spec = load(SCHEMA)
        self.defs = self.spec['$defs']

    def fixture(self, name):
        return load(DIRECTORY / 'fixtures/valid' / name)

    def check(self, message):
        validate(message, self.spec, self.spec)

    def test_upstream_pin(self):
        message = self.fixture('offline-grant.json')
        self.check(message)
        message['body']['aegivela_revision'] = 'latest'
        with self.assertRaises(SchemaError):
            self.check(message)

    def test_opaque_relay(self):
        message = self.fixture('offline-grant.json')
        raw = message['body']['signed_offline_grant']
        self.assertEqual(raw, 'eyJhIjogMX0.eyJiIjoyfQ.c2ln')
        self.assertEqual(hashlib.sha256(raw.encode('ascii')).hexdigest(),
                         message['body']['artifact_digest'])
        self.check(message)
        for field in ('private_key', 'scope', 'budget', 'unsigned_grant'):
            candidate = copy.deepcopy(message)
            candidate['body'][field] = 'caller-controlled'
            with self.subTest(field=field), self.assertRaises(SchemaError):
                self.check(candidate)

    def test_digest_domains(self):
        message = self.fixture('local-admission-batch.json')
        self.check(message)
        for name in ('artifact_digest', 'authority_binding_digest',
                     'authority_payload_digest', 'operation_binding_digest',
                     'operation_payload_digest'):
            candidate = copy.deepcopy(message)
            fact = candidate['body']['facts'][0]
            value = fact[name]
            fact[name] = value.removeprefix('sha256:') if value.startswith('sha256:') else 'sha256:' + value
            with self.subTest(field=name), self.assertRaises(SchemaError):
                self.check(candidate)

    def test_scope_binding(self):
        message = self.fixture('local-admission-batch.json')
        for field in ('tenant_id', 'workspace_id', 'site_id', 'satellite_id',
                      'source_epoch', 'coordinator_epoch'):
            candidate = copy.deepcopy(message)
            del candidate[field]
            with self.subTest(field=field), self.assertRaises(SchemaError):
                self.check(candidate)
        for field in ('namespace', 'offline_grant_id', 'admission_id',
                      'operation_key', 'attempt_id'):
            candidate = copy.deepcopy(message)
            del candidate['body']['facts'][0][field]
            with self.subTest(field=field), self.assertRaises(SchemaError):
                self.check(candidate)

    def test_manifest_ids(self):
        mapping = load(DIRECTORY / 'upstream.json')
        self.assertEqual(mapping['aegivela_revision'], REVISION)
        self.assertEqual(mapping['aegivela_contract'], 'aegivela.io/offline-admission/v1.0')
        cases = mapping['consumer_cases']
        self.assertEqual(set(cases), CASE_IDS)
        self.assertEqual(len(cases), len(CASE_IDS))
        self.assertEqual(mapping['execution_status'], 'not-executed')

    def test_capability_ack(self):
        message = self.fixture('hello-ack.json')
        self.check(message)
        for capabilities in ([], ['submit'], ['offline_admission_v1', 'offline_admission_v1']):
            candidate = copy.deepcopy(message)
            candidate['body']['capabilities'] = capabilities
            with self.subTest(capabilities=capabilities), self.assertRaises(SchemaError):
                self.check(candidate)

    def test_unreviewed_messages_do_not_widen_legacy_v1(self):
        legacy = load(DIRECTORY.parent / 'satellite.schema.json')
        message = self.fixture('offline-grant.json')
        with self.assertRaises(SchemaError):
            validate(message, legacy['$defs']['EnvelopeBase'], legacy)

    def test_message_body_pairing(self):
        message = self.fixture('offline-grant.json')
        for kind in ('offline_grant_ack', 'local_admission_batch', 'local_admission_ack', 'submit'):
            candidate = copy.deepcopy(message)
            candidate['message_type'] = kind
            with self.subTest(kind=kind), self.assertRaises(SchemaError):
                self.check(candidate)

    def test_batch_bounds_and_uncertainty(self):
        message = self.fixture('local-admission-batch.json')
        for state in ('admitted', 'start_intent', 'terminal', 'unknown', 'quarantined'):
            message['body']['facts'][0]['state'] = state
            self.check(message)
        message['body']['facts'][0]['state'] = 'effect_verified'
        with self.assertRaises(SchemaError):
            self.check(message)
        message = self.fixture('local-admission-batch.json')
        for facts in ([], message['body']['facts'] * 101):
            candidate = copy.deepcopy(message)
            candidate['body']['facts'] = facts
            with self.subTest(count=len(facts)), self.assertRaises(SchemaError):
                self.check(candidate)

    def test_exact_integer_charges(self):
        message = self.fixture('local-admission-batch.json')
        for value in (-1, 0.5, 9007199254740992, True):
            candidate = copy.deepcopy(message)
            candidate['body']['facts'][0]['reserved_resources']['duration_ms'] = value
            with self.subTest(value=value), self.assertRaises(SchemaError):
                self.check(candidate)

    def test_fact_ack_requires_original_source_epoch(self):
        for name in ('local-admission-batch.json', 'local-admission-ack.json'):
            message = self.fixture(name)
            self.check(message)
            message['body'].pop('admission_source_epoch', None)
            with self.subTest(message=name), self.assertRaises(SchemaError):
                self.check(message)

    def test_fact_ack_requires_namespace(self):
        for name in ('local-admission-batch.json', 'local-admission-ack.json'):
            message = self.fixture(name)
            self.check(message)
            message['body'].pop('namespace', None)
            with self.subTest(message=name), self.assertRaises(SchemaError):
                self.check(message)

    def test_capability_ack_requires_boolean_acceptance(self):
        message = self.fixture('hello-ack.json')
        message['body']['accepted'] = 1
        with self.assertRaises(SchemaError):
            self.check(message)

    def test_valid_fixtures(self):
        paths = sorted((DIRECTORY / 'fixtures/valid').glob('*.json'))
        self.assertEqual(len(paths), 5)
        for path in paths:
            with self.subTest(path=path.name):
                self.check(load(path))

    def test_invalid_fixtures(self):
        paths = sorted((DIRECTORY / 'fixtures/invalid').glob('*.json'))
        self.assertTrue(paths)
        for path in paths:
            with self.subTest(path=path.name), self.assertRaises(SchemaError):
                self.check(load(path))


class SchemaValidationBoundaryTests(unittest.TestCase):
    def test_one_of_rejects_no_matching_message(self):
        schema = {'oneOf': [{'const': 'one'}, {'const': 'two'}]}
        with self.assertRaises(SchemaError):
            validate('three', schema, schema)

    def test_one_of_rejects_ambiguous_message(self):
        schema = {'oneOf': [{'type': 'integer'}, {'type': 'number'}]}
        with self.assertRaises(SchemaError):
            validate(1, schema, schema)

    def test_maximum_batch_length_is_enforced(self):
        schema = {'type': 'array', 'maxItems': 1}
        with self.assertRaises(SchemaError):
            validate([1, 2], schema, schema)

    def test_duplicate_capabilities_are_rejected(self):
        schema = {'type': 'array', 'uniqueItems': True}
        with self.assertRaises(SchemaError):
            validate(['offline_admission_v1'] * 2, schema, schema)


if __name__ == '__main__':
    unittest.main()
