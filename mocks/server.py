import json
import logging
import os
import re
import signal
import sys
import threading
import urllib.request
from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from profiles import InvalidINN, account, credit_report, fault, registry

ROUTES = {
    re.compile(r"^/fns/v1/companies/([^/]+)$"): ("fns", registry),
    re.compile(r"^/bki/v1/reports/([^/]+)$"): ("bki", credit_report),
    re.compile(r"^/bank/v1/accounts/([^/]+)$"): ("bank", account),
}


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/health":
            return self._send(HTTPStatus.OK, {"status": "SERVING"})

        for pattern, (source, build) in ROUTES.items():
            match = pattern.match(self.path)
            if not match:
                continue
            inn = match.group(1)
            if fault(inn, source) == "UNAVAILABLE":
                return self._send(HTTPStatus.SERVICE_UNAVAILABLE, {"error": f"{source} недоступен"})
            try:
                return self._send(HTTPStatus.OK, build(inn))
            except InvalidINN:
                return self._send(HTTPStatus.BAD_REQUEST, {"error": "ИНН должен состоять из 10 или 12 цифр"})

        self._send(HTTPStatus.NOT_FOUND, {"error": "Not found"})

    def _send(self, status: HTTPStatus, body: dict):
        payload = json.dumps(body, ensure_ascii=False).encode()
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def log_message(self, format, *args):
        logging.info("%s %s", self.command, args[1] if len(args) > 1 else "")


def main():
    port = int(os.getenv("PORT", "8090"))

    if len(sys.argv) > 1 and sys.argv[1] == "healthcheck":
        with urllib.request.urlopen(f"http://127.0.0.1:{port}/health", timeout=2) as reply:
            raise SystemExit(0 if reply.status == 200 else 1)

    logging.basicConfig(level=logging.INFO, format="%(message)s")
    server = ThreadingHTTPServer(("0.0.0.0", port), Handler)
    signal.signal(signal.SIGTERM, lambda *_: threading.Thread(target=server.shutdown).start())
    logging.info("external mocks listening on %s", port)
    server.serve_forever()


if __name__ == "__main__":
    main()
