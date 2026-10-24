import ast
import os
from pathlib import Path
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[1]


class RemoveDocstrings(ast.NodeTransformer):
    def clean(self, node):
        self.generic_visit(node)

        if getattr(node, "body", None) and isinstance(node.body[0], ast.Expr) and isinstance(node.body[0].value, ast.Constant) and isinstance(node.body[0].value.value, str):
            node.body.pop(0)

            if not node.body:
                node.body.append(ast.Pass())

        return node
    visit_Module = clean
    visit_ClassDef = clean
    visit_FunctionDef = clean
    visit_AsyncFunctionDef = clean


def strip():
    files = list((ROOT / "contracts").rglob("*.go"))
    subprocess.run(["go", "run", str(ROOT / "scripts/strip-comments.go"), "--", *map(str, files)], cwd=ROOT, check=True)

    for path in (ROOT / "scoring/generated").rglob("*.py"):
        tree = RemoveDocstrings().visit(ast.parse(path.read_text()))
        path.write_text(ast.unparse(ast.fix_missing_locations(tree)) + "\n")

if __name__ == "__main__":
    if "--strip-only" in sys.argv:
        strip()
    else:
        with tempfile.TemporaryDirectory() as temp:
            env = {**os.environ, "GOBIN": temp, "GOWORK": "off", "PATH": temp + os.pathsep + os.environ["PATH"]}

            for plugin in ("google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1"):
                subprocess.run(["go", "install", plugin], env=env, check=True)

            subprocess.run([sys.executable, "-m", "grpc_tools.protoc", "-I", str(ROOT / "api/proto"), "--go_out=" + str(ROOT / "contracts"), "--go_opt=module=robust/contracts", "--go-grpc_out=" + str(ROOT / "contracts"), "--go-grpc_opt=module=robust/contracts", "--python_out=" + str(ROOT / "scoring/generated"), "--grpc_python_out=" + str(ROOT / "scoring/generated"), "credit/v1/credit.proto"], env=env, check=True)

        strip()
