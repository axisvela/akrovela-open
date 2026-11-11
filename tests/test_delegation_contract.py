import json
from pathlib import Path
import unittest

from jsonschema import Draft202012Validator, FormatChecker

ROOT = Path(__file__).resolve().parents[1]
SCHEMA = ROOT / 'contracts/delegation/v1/delegated-operation.schema.json'
FIXTURES = ROOT / 'contracts/delegation/v1/fixtures'
OPERATIONS = ROOT / 'contracts/operations/v1/operations.openapi.json'

REQUIRED = [
    'schema_version', 'contract_mode', 'delegation_id', 'operation_key',
    'operation_type', 'scope', 'mission_ref', 'assignment_ref', 'lease_ref',
    'fencing_token', 'authority_ref', 'safety_envelope_ref', 'approval_ref',
    'expected_effect', 'target', 'payload', 'observation_refs', 'validity',
    'resource_limits', 'idempotency_retention_seconds',
]

class DelegationContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.schema = json.loads(SCHEMA.read_text())
        cls.validator = Draft202012Validator(cls.schema, format_checker=FormatChecker())

    def test_identity_and_mode(self):
        self.assertEqual('https://schemas.axisvela.com/akrovela-delegated-operation.v1.schema.json', self.schema['$id'])
        self.assertEqual('akrovela.delegated-operation/v1', self.schema['properties']['schema_version']['const'])
        self.assertEqual(['live', 'fixture'], self.schema['properties']['contract_mode']['enum'])
        self.assertEqual(set(REQUIRED), set(self.schema['required']))

    def test_defs_match_operations_components(self):
        operations = json.loads(OPERATIONS.read_text())['components']['schemas']
        objects = ['Scope', 'Target', 'Validity', 'ResourceLimits', 'NavigateToPosePayload', 'ObservationReference', 'ObservationSubject']
        strings = ['ObservationId', 'Digest']
        for name in objects:
            with self.subTest(schema=name):
                self.assertIn(name, self.schema['$defs'])
                self.assertEqual(set(operations[name]['required']), set(self.schema['$defs'][name]['required']))
                self.assertEqual(set(operations[name]['properties']), set(self.schema['$defs'][name]['properties']))
        for name in strings:
            with self.subTest(schema=name):
                self.assertIn(name, self.schema['$defs'])
                self.assertEqual(operations[name]['type'], self.schema['$defs'][name]['type'])
                self.assertEqual(operations[name].get('pattern'), self.schema['$defs'][name].get('pattern'))

    def test_local_refs_only(self):
        def walk(value):
            if isinstance(value, dict):
                ref = value.get('$ref')
                if isinstance(ref, str):
                    self.assertTrue(ref.startswith('#/'), ref)
                for child in value.values():
                    walk(child)
            elif isinstance(value, list):
                for child in value:
                    walk(child)
        walk(self.schema)

    def test_valid_fixtures(self):
        paths = sorted((FIXTURES / 'valid').glob('*.json'))
        self.assertGreaterEqual(len(paths), 2)
        for path in paths:
            with self.subTest(fixture=path.name):
                self.assertEqual([], list(self.validator.iter_errors(json.loads(path.read_text()))))

    def test_invalid_fixtures(self):
        paths = sorted((FIXTURES / 'invalid').glob('*.json'))
        self.assertGreaterEqual(len(paths), 6)
        for path in paths:
            with self.subTest(fixture=path.name):
                self.assertNotEqual([], list(self.validator.iter_errors(json.loads(path.read_text()))))

    def test_review_record_exists(self):
        text = (ROOT / 'contracts/delegation/v1/review.md').read_text()
        self.assertIn('## Ownership', text)
        self.assertIn('## Negative and Fault Cases', text)
        self.assertIn('not external', text.lower())

if __name__ == '__main__':
    unittest.main()
