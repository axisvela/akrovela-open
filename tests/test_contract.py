import json
from pathlib import Path
import re
import unittest

ROOT = Path(__file__).resolve().parents[1]


class ContractTests(unittest.TestCase):
    def test_bootstrap_contract(self):
        spec = json.loads((ROOT / 'contracts/bootstrap.openapi.json').read_text())
        version = (ROOT / 'VERSION').read_text().strip()
        self.assertRegex(version, r'^\d+\.\d+\.\d+$')
        self.assertEqual(spec['info']['version'], version)
        self.assertEqual(set(spec['paths']), {'/healthz', '/readyz', '/version'})
        self.assertIn('503', spec['paths']['/readyz']['get']['responses'])
        schemas = spec['components']['schemas']
        self.assertEqual(schemas['Version']['properties']['product']['const'], 'AKROVELA')
        self.assertTrue(re.fullmatch(schemas['Version']['properties']['version']['pattern'], version))

    def test_local_references_resolve(self):
        spec = json.loads((ROOT / 'contracts/bootstrap.openapi.json').read_text())

        def visit(value):
            if isinstance(value, dict):
                if '$ref' in value:
                    self.assertTrue(value['$ref'].startswith('#/'))
                    target = spec
                    for part in value['$ref'][2:].split('/'):
                        target = target[part]
                for child in value.values():
                    visit(child)
            elif isinstance(value, list):
                for child in value:
                    visit(child)

        visit(spec)


if __name__ == '__main__':
    unittest.main()
