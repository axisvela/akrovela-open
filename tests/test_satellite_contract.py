import sys
from pathlib import Path
import unittest

sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_operations_contract import SchemaError, load, validate

ROOT = Path(__file__).resolve().parents[1]
SCHEMA = ROOT / 'contracts/satellite/v1/satellite.schema.json'
FIXTURES = ROOT / 'contracts/satellite/v1/fixtures'


class SatelliteContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.spec = load(SCHEMA)
        cls.defs = cls.spec['$defs']

    def test_contract_identity(self):
        self.assertEqual(self.spec['x-contract-version'], 'v1')
        self.assertEqual(self.spec['x-contract-status'], 'reviewed')
        self.assertEqual(self.spec['x-protocol-version'], 'akrovela.satellite/v1')

    def test_message_set(self):
        for name in (
            'HelloMessage', 'SubmitMessage', 'SubmitResultMessage',
            'CancelMessage', 'CancelResultMessage', 'ReceiptBatchMessage',
            'HeartbeatMessage', 'EnvelopeBase', 'SubmitBody', 'Target',
            'Validity', 'ResourceLimits',
        ):
            self.assertIn(name, self.defs)

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
