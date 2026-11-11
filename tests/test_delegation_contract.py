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

DEFS = [
    'OperationKey', 'AttemptId', 'ObservationId', 'Digest', 'Scope', 'Target',
    'Validity', 'ResourceLimits', 'ObservationSubject', 'ObservationReference',
    'NavigateToPosePayload',
]

INVALID_REASONS = {
    'delegated-operation-missing-delegation-id': ('required', ['delegation_id']),
    'delegated-operation-wrong-schema-version': ('const', ['schema_version']),
    'delegated-operation-bad-contract-mode': ('enum', ['contract_mode']),
    'delegated-operation-zero-fence': ('minimum', ['fencing_token']),
    'delegated-operation-missing-mission-ref': ('required', ['mission_ref']),
    'delegated-operation-missing-expected-effect': ('required', ['expected_effect']),
}


def rewrite_component_refs(value):
    if isinstance(value, dict):
        rewritten = {}
        for key, child in value.items():
            prefix = '#/components/schemas/'
            if key == '$ref' and isinstance(child, str) and child.startswith(prefix):
                rewritten[key] = '#/$defs/' + child[len(prefix):]
            else:
                rewritten[key] = rewrite_component_refs(child)
        return rewritten
    if isinstance(value, list):
        return [rewrite_component_refs(item) for item in value]
    return value


def error_path(error):
    path = list(error.absolute_path)
    if error.validator == 'required':
        path.append(error.message.split("'")[1])
    return path


class DelegationContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.schema = json.loads(SCHEMA.read_text(encoding='utf-8'))
        cls.validator = Draft202012Validator(cls.schema, format_checker=FormatChecker())

    def test_identity_and_mode(self):
        self.assertEqual('https://schemas.axisvela.com/akrovela-delegated-operation.v1.schema.json', self.schema['$id'])
        self.assertEqual('akrovela.delegated-operation/v1', self.schema['properties']['schema_version']['const'])
        self.assertEqual(['live', 'fixture'], self.schema['properties']['contract_mode']['enum'])
        self.assertEqual(set(REQUIRED), set(self.schema['required']))

    def test_defs_match_operations_components(self):
        operations = json.loads(OPERATIONS.read_text(encoding='utf-8'))['components']['schemas']
        for name in DEFS:
            with self.subTest(schema=name):
                self.assertIn(name, self.schema['$defs'])
                self.assertIn(name, operations)
                self.assertEqual(
                    rewrite_component_refs(operations[name]),
                    self.schema['$defs'][name],
                )

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
                document = json.loads(path.read_text(encoding='utf-8'))
                self.assertEqual([], list(self.validator.iter_errors(document)))

    def test_invalid_fixtures(self):
        paths = sorted((FIXTURES / 'invalid').glob('*.json'))
        self.assertEqual(set(INVALID_REASONS), {path.stem for path in paths})
        for path in paths:
            with self.subTest(fixture=path.name):
                document = json.loads(path.read_text(encoding='utf-8'))
                errors = list(self.validator.iter_errors(document))
                self.assertNotEqual([], errors)
                first = errors[0]
                self.assertEqual(
                    (first.validator, error_path(first)),
                    INVALID_REASONS[path.stem],
                )

    def test_review_record_exists(self):
        text = (ROOT / 'contracts/delegation/v1/review.md').read_text(encoding='utf-8')
        self.assertIn('## Ownership', text)
        self.assertIn('## Negative and Fault Cases', text)
        self.assertIn('not external', text.lower())

if __name__ == '__main__':
    unittest.main()
