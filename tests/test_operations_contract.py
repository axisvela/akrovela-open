import hashlib
import json
from pathlib import Path
import re
import unittest

ROOT = Path(__file__).resolve().parents[1]
CONTRACT = ROOT / 'contracts/operations/v1/operations.openapi.json'
FIXTURES = ROOT / 'contracts/operations/v1/fixtures'

BINDING_FIELDS = [
    'scope', 'target', 'mission_ref', 'assignment_ref', 'lease_ref',
    'fencing_token', 'safety_envelope_ref', 'authority_ref', 'approval_ref',
    'validity', 'resource_limits', 'observation_refs',
    'idempotency_retention_seconds',
]


class SchemaError(AssertionError):
    pass


def load(path):
    return json.loads(path.read_text())


def canonical(value):
    if isinstance(value, dict):
        return '{' + ','.join(
            json.dumps(key, ensure_ascii=False) + ':' + canonical(value[key])
            for key in sorted(value)
        ) + '}'
    if isinstance(value, list):
        return '[' + ','.join(canonical(item) for item in value) + ']'
    if isinstance(value, bool):
        return 'true' if value else 'false'
    if value is None:
        return 'null'
    if isinstance(value, int):
        return str(value)
    if isinstance(value, float):
        if value == int(value) and abs(value) < 1e21:
            return str(int(value))
        return repr(value)
    if isinstance(value, str):
        return json.dumps(value, ensure_ascii=False)
    raise TypeError(type(value))


def binding_digest(document, fields):
    binding = {name: document[name] for name in fields}
    return 'sha256:' + hashlib.sha256(canonical(binding).encode('utf-8')).hexdigest()


def resolve(root, schema):
    while '$ref' in schema:
        ref = schema['$ref']
        if not ref.startswith('#/'):
            raise SchemaError(f'external ref {ref}')
        target = root
        for part in ref[2:].split('/'):
            target = target[part]
        schema = target
    return schema


def matches(value, schema, root):
    try:
        validate(value, schema, root)
        return True
    except SchemaError:
        return False


def validate(value, schema, root):
    schema = resolve(root, schema)
    if 'allOf' in schema:
        for child in schema['allOf']:
            validate(value, child, root)
    if 'if' in schema and matches(value, schema['if'], root):
        validate(value, schema.get('then', {}), root)
    if 'const' in schema and value != schema['const']:
        raise SchemaError(f'expected const {schema["const"]!r}, got {value!r}')
    if 'enum' in schema and value not in schema['enum']:
        raise SchemaError(f'expected one of {schema["enum"]!r}, got {value!r}')
    expected = schema.get('type')
    if expected is not None:
        types = expected if isinstance(expected, list) else [expected]
        valid_type = False
        for item in types:
            if item == 'object' and isinstance(value, dict):
                valid_type = True
            elif item == 'array' and isinstance(value, list):
                valid_type = True
            elif item == 'string' and isinstance(value, str):
                valid_type = True
            elif item == 'integer' and isinstance(value, int) and not isinstance(value, bool):
                valid_type = True
            elif item == 'number' and isinstance(value, (int, float)) and not isinstance(value, bool):
                valid_type = True
            elif item == 'boolean' and isinstance(value, bool):
                valid_type = True
            elif item == 'null' and value is None:
                valid_type = True
        if not valid_type:
            raise SchemaError(f'expected type {expected!r}, got {type(value).__name__}')
    if isinstance(value, dict):
        for key in schema.get('required', []):
            if key not in value:
                raise SchemaError(f'missing required {key!r}')
        properties = schema.get('properties', {})
        if schema.get('additionalProperties') is False:
            extra = set(value) - set(properties)
            if extra:
                raise SchemaError(f'unexpected properties {sorted(extra)}')
        for key, child in value.items():
            if key in properties:
                validate(child, properties[key], root)
    if isinstance(value, list):
        if len(value) < schema.get('minItems', 0):
            raise SchemaError('too few items')
        if 'items' in schema:
            for item in value:
                validate(item, schema['items'], root)
    if isinstance(value, str):
        if len(value) < schema.get('minLength', 0):
            raise SchemaError('string too short')
        if 'pattern' in schema and not re.search(schema['pattern'], value):
            raise SchemaError(f'{value!r} does not match {schema["pattern"]!r}')
        if schema.get('format') == 'date-time' and not re.fullmatch(
            r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z', value
        ):
            raise SchemaError(f'{value!r} is not a UTC date-time')
    if isinstance(value, (int, float)) and not isinstance(value, bool):
        if 'minimum' in schema and value < schema['minimum']:
            raise SchemaError('below minimum')
        if 'maximum' in schema and value > schema['maximum']:
            raise SchemaError('above maximum')


class OperationsContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.spec = load(CONTRACT)
        cls.schemas = cls.spec['components']['schemas']

    def test_contract_identity_and_version(self):
        self.assertEqual(self.spec['openapi'], '3.1.0')
        self.assertEqual(self.spec['info']['version'], (ROOT / 'VERSION').read_text().strip())
        self.assertEqual(self.spec['info']['x-contract-version'], 'v1')
        self.assertEqual(self.spec['info']['x-contract-status'], 'reviewed')
        self.assertIn('bearerAuth', self.spec['components']['securitySchemes'])

    def test_paths_and_operations(self):
        expected = {
            ('/v1/operations', 'post'),
            ('/v1/operations/{operation_key}', 'get'),
            ('/v1/operations/{operation_key}/cancel', 'post'),
            ('/v1/operations/{operation_key}/receipts', 'get'),
            ('/v1/receipts/{receipt_id}', 'get'),
        }
        found = {
            (path, method)
            for path, item in self.spec['paths'].items()
            for method in item
            if method in {'get', 'post'}
        }
        self.assertEqual(found, expected)

    def test_local_references_resolve(self):
        def visit(value):
            if isinstance(value, dict):
                if '$ref' in value:
                    self.assertTrue(value['$ref'].startswith('#/'))
                    target = self.spec
                    for part in value['$ref'][2:].split('/'):
                        target = target[part]
                for child in value.values():
                    visit(child)
            elif isinstance(value, list):
                for child in value:
                    visit(child)

        visit(self.spec)

    def test_required_schemas(self):
        for name in (
            'Scope', 'Target', 'Validity', 'ResourceLimits',
            'ObservationReference', 'ObservationSubject', 'NavigateToPosePayload',
            'SubmitOperationRequest', 'OperationView', 'AttemptView',
            'CancelRequest', 'CancelResponse', 'Receipt', 'ReceiptList',
            'ExecutorFacts', 'VerifierResult', 'Error',
        ):
            self.assertIn(name, self.schemas)

    def test_state_axes_are_separate(self):
        operation = self.schemas['OperationView']
        required = set(operation['required'])
        for name in (
            'admission_status', 'execution_status',
            'cancellation_status', 'effect_status',
        ):
            self.assertIn(name, required)
        self.assertEqual(
            set(operation['properties']['admission_status']['enum']),
            {'accepted', 'rejected'},
        )
        self.assertIn('confirmed_no_effect', operation['properties']['execution_status']['enum'])
        self.assertIn('stop_confirmed', operation['properties']['cancellation_status']['enum'])
        self.assertIn('verified', operation['properties']['effect_status']['enum'])

    def test_binding_digest_covers_all_immutable_effective_fields(self):
        binding = self.schemas['SubmitOperationRequest']['properties']['binding_digest']
        self.assertEqual(binding.get('x-binding-fields'), BINDING_FIELDS)

    def test_valid_submit_binding_digest_matches_fixture(self):
        document = load(FIXTURES / 'valid' / 'submit-navigate.json')
        self.assertEqual(document['binding_digest'], binding_digest(document, BINDING_FIELDS))

    def test_submit_bindings(self):
        request = self.schemas['SubmitOperationRequest']
        for name in (
            'operation_key', 'payload_digest', 'binding_digest', 'target',
            'fencing_token', 'validity', 'resource_limits', 'authority_ref',
            'safety_envelope_ref',
        ):
            self.assertIn(name, request['required'])
        pattern = resolve(self.spec, request['properties']['operation_key'])['pattern']
        self.assertIsNone(re.search(pattern, 'op:tenant/acme'))
        self.assertIsNotNone(re.search(pattern, 'op:0192f0f2-7d6f-7c2a-9d1e-2b6f0f7a9c31'))
        self.assertEqual(
            self.schemas['Validity']['properties']['clock_domain']['const'],
            'authority-utc',
        )

    def test_cancel_stop_requires_source(self):
        cancel = self.schemas['CancelResponse']
        self.assertIn('allOf', cancel)
        self.assertIn('stop_confirmed', cancel['required'])
        self.assertIn('confirmation_source', cancel['required'])
        self.assertNotIn('not_found', cancel['properties']['cancel_status']['enum'])

    def test_security_error_codes(self):
        codes = set(
            self.schemas['Error']['properties']['error']['properties']['code']['enum']
        )
        for code in (
            'payload_digest_mismatch', 'binding_digest_mismatch',
            'authority_revoked', 'approval_consumed', 'approval_expired',
            'stale_fence', 'stale_node_boot_epoch', 'forbidden_scope',
        ):
            self.assertIn(code, codes)

    def test_receipt_types_and_conditionals(self):
        receipt = self.schemas['Receipt']
        self.assertIn('receipt_type', receipt['required'])
        self.assertIn('receipt_type', receipt['properties'])
        self.assertIn('confirmation_source', receipt['properties'])
        self.assertIn('allOf', receipt)

    def test_valid_fixtures(self):
        schema_names = {
            'submit-navigate': 'SubmitOperationRequest',
            'operation-view': 'OperationView',
            'cancel-request': 'CancelRequest',
            'cancel-response': 'CancelResponse',
            'receipt': 'Receipt',
            'receipt-admission-rejected': 'Receipt',
        }
        for path in sorted((FIXTURES / 'valid').glob('*.json')):
            validate(load(path), self.schemas[schema_names[path.stem]], self.spec)

    def test_invalid_fixtures_rejected(self):
        for path in sorted((FIXTURES / 'invalid').glob('*.json')):
            document = load(path)
            schema_name = document.pop('_schema')
            with self.assertRaises(SchemaError, msg=path.name):
                validate(document, self.schemas[schema_name], self.spec)


if __name__ == '__main__':
    unittest.main()
