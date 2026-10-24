from pathlib import Path
import re
import warnings
import yaml
from jsonschema import Draft4Validator

with warnings.catch_warnings():
    warnings.simplefilter("ignore", DeprecationWarning)
    from jsonschema import RefResolver

from openapi_spec_validator import validate_spec

SPEC = yaml.safe_load((Path(__file__).resolve().parents[1] / "api/public/openapi.yaml").read_text())
validate_spec(SPEC)


def check_response(method, path, status, body):
    path = path.split("?", 1)[0]

    for template, operations in SPEC["paths"].items():
        pattern = re.sub(r"\{[^}]+\}", "[^/]+", template)

        if re.fullmatch(pattern, path) and method.lower() in operations:
            operation = operations[method.lower()]
            response = operation["responses"].get(str(status))
            assert response is not None, (method, path, "undocumented status", status)
            schema = response.get("content", {}).get("application/json", {}).get("schema")

            if schema:
                Draft4Validator(schema, resolver=RefResolver.from_schema(SPEC)).validate(body)

            return

    raise AssertionError(("undocumented operation", method, path))

if __name__ == "__main__":
    print("OpenAPI is valid")
