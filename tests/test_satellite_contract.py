import sys
from pathlib import Path
import unittest

sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_operations_contract import SchemaError, load, validate

ROOT = Path(__file__).resolve().parents[1]
SCHEMA = ROOT / 'contracts/satellite/v1/satellite.schema.json'
FIXTURES = ROOT / 'contracts/satellite/v1/fixtures'
VECTORS = ROOT / 'contracts/satellite/v1/vectors/semantic-vectors.json'


class SatelliteContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.spec = load(SCHEMA)
        cls.defs = cls.spec['$defs']

    def valid_body(self, document):
        schema_name = document['_schema']
        body = document['body']
        for part in self.defs[schema_name].get('allOf', ()):
            ref = part.get('properties', {}).get('body', {}).get('$ref')
            if ref is None:
                continue
            try:
                validate(body, {'$ref': ref}, self.spec)
            except SchemaError:
                return False
            if (body.get('execution_model') == 'repeat_attempt'
                    and body.get('contract_mode') != 'fixture'):
                return False
            return True
        raise SchemaError(f'{schema_name} has no body schema')

    def test_contract_identity(self):
        self.assertEqual(self.spec['x-contract-version'], 'v1')
        self.assertEqual(self.spec['x-contract-status'], 'reviewed')
        self.assertEqual(self.spec['x-protocol-version'], 'akrovela.satellite/v1')

    def test_message_set(self):
        for name in (
            'HelloMessage', 'HelloAckMessage', 'SubmitMessage', 'SubmitResultMessage',
            'CancelMessage', 'CancelResultMessage', 'ReceiptBatchMessage',
            'ReceiptAckMessage', 'HeartbeatMessage', 'StatusMessage',
            'StatusResultMessage', 'ResumeMessage', 'ResumeResultMessage',
            'PreReleaseMessage', 'PreReleaseAckMessage', 'ConsumptionBatchMessage',
            'ConsumptionAckMessage', 'PreReleaseBody', 'ConsumptionBatchBody',
            'EnvelopeBase', 'SubmitBody', 'Target', 'Validity', 'ResourceLimits',
            'ErrorCode',
        ):
            self.assertIn(name, self.defs)

    def test_pre_release_requires_budget(self):
        body = self.defs['PreReleaseBody']
        for name in ('allocation_id', 'count_budget', 'resource_budget',
                     'binding_digest', 'dispatch_not_after', 'contract_mode'):
            self.assertIn(name, body['required'])
        self.assertEqual(body['properties']['contract_mode']['const'], 'fixture')

    def test_pre_release_multi_fields(self):
        body = load(FIXTURES / 'valid' / 'pre-release-multi.json')['body']
        self.assertEqual(body['execution_model'], 'satellite_attempts')
        self.assertEqual(body['resource_budget_mode'], 'aggregate')
        self.assertEqual(body['count_budget'], 3)

    def test_repeat_release_requires_fixture(self):
        valid = load(FIXTURES / 'valid' / 'pre-release-repeat.json')['body']
        self.assertEqual(valid['execution_model'], 'repeat_attempt')
        invalid = load(FIXTURES / 'invalid' / 'pre-release-repeat-live.json')
        self.assertFalse(self.valid_body(invalid))

    def test_pre_release_count_budget_maximum(self):
        body = self.defs['PreReleaseBody']
        self.assertEqual(body['properties']['count_budget']['maximum'], 10)
        document = load(FIXTURES / 'invalid' / 'pre-release-count-budget-eleven.json')
        schema_name = document.pop('_schema')
        with self.assertRaises(SchemaError):
            validate(document, self.defs[schema_name], self.spec)

    def test_receipt_evidence_fields(self):
        document = load(FIXTURES / 'valid' / 'receipt-batch-evidence.json')
        schema_name = document.pop('_schema')
        validate(document, self.defs[schema_name], self.spec)
        receipt = document['body']['receipts'][0]
        self.assertEqual(receipt['receipt_type'], 'evidence')
        self.assertEqual(receipt['execution_index'], 2)
        self.assertEqual(receipt['repeat_of_attempt_id'], receipt['attempt_id'])

    def test_consumption_indexed_fields(self):
        body = load(FIXTURES / 'valid' / 'consumption-batch-indexed.json')['body']
        self.assertEqual(body['consumptions'][0]['execution_index'], 2)
        self.assertEqual(body['consumptions'][0]['resource_usage']['duration_ms'], 60000)

    def test_consumption_ack_results(self):
        body = load(FIXTURES / 'valid' / 'consumption-ack-results.json')['body']
        self.assertEqual(body['results'][0]['status'], 'accepted')
        self.assertEqual(body['results'][1]['reason'], 'execution_index_conflict')

    def test_error_code_set(self):
        codes = set(self.defs['ErrorCode']['enum'])
        for code in (
            'unsupported_version', 'forbidden_scope', 'expired_dispatch',
            'stale_coordinator_epoch', 'stale_source_epoch', 'binding_conflict',
            'queue_full', 'frame_too_large', 'reconciliation_required',
            'action_in_progress', 'cancel_not_possible', 'unknown_outcome',
        ):
            self.assertIn(code, codes)

    def test_envelope_requires_identity_and_epoch(self):
        base = self.defs['EnvelopeBase']
        for name in (
            'protocol_version', 'message_type', 'message_id', 'session_id',
            'sequence', 'tenant_id', 'workspace_id', 'site_id', 'satellite_id',
            'source_epoch', 'coordinator_id', 'coordinator_epoch', 'sent_at', 'body',
        ):
            self.assertIn(name, base['required'])
        self.assertIs(base['additionalProperties'], False)

    def test_submit_requires_fence_and_dispatch(self):
        body = self.defs['SubmitBody']
        for name in ('operation_key', 'attempt_id', 'binding_digest', 'target',
                     'fencing_token', 'validity', 'dispatch_not_after',
                     'contract_mode'):
            self.assertIn(name, body['required'])
        self.assertEqual(body['properties']['contract_mode']['const'], 'fixture')

    def test_semantic_vectors(self):
        document = load(VECTORS)
        self.assertEqual(document['protocol_version'], 'akrovela.satellite/v1')
        codes = set(self.defs['ErrorCode']['enum'])
        names = set()
        for vector in document['vectors']:
            self.assertIn('name', vector)
            names.add(vector['name'])
            self.assertIn('message', vector)
            self.assertIn('now', vector)
            self.assertIn('clock_uncertainty_ms', vector)
            code = vector['expect_code']
            self.assertTrue(code == '' or code in codes, code)
        self.assertIn('valid', names)
        self.assertIn('expired-dispatch', names)
        self.assertIn('binding-digest-mismatch', names)
        valid = next(vector for vector in document['vectors'] if vector['name'] == 'valid')
        validate(valid['message'], self.defs['SubmitMessage'], self.spec)

    def test_valid_fixtures(self):
        for path in sorted((FIXTURES / 'valid').glob('*.json')):
            document = load(path)
            schema_name = document.pop('_schema')
            validate(document, self.defs[schema_name], self.spec)

    def test_invalid_fixtures_rejected(self):
        paths = sorted((FIXTURES / 'invalid').glob('*.json'))
        self.assertTrue(paths, 'no invalid fixtures')
        for path in paths:
            document = load(path)
            schema_name = document.pop('_schema')
            with self.assertRaises(SchemaError, msg=path.name):
                validate(document, self.defs[schema_name], self.spec)


if __name__ == '__main__':
    unittest.main()
