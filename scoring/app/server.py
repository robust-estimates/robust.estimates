import base64
import hashlib
import hmac
import json
import logging
import os
import signal
import sys
import time
from concurrent import futures

import grpc
from credit.v1 import credit_pb2 as pb
from credit.v1 import credit_pb2_grpc as rpc
from grpc_health.v1 import health, health_pb2, health_pb2_grpc
from app.model import RiskModel


def authenticate(context):
    raw = dict(context.invocation_metadata()).get("authorization", "")

    try:
        token = raw.removeprefix("Bearer ")
        header, payload, signature = token.split(".")
        decode = lambda value: base64.urlsafe_b64decode(value + "=" * (-len(value) % 4))
        expected = hmac.new(os.environ["JWT_SECRET"].encode(), f"{header}.{payload}".encode(), hashlib.sha256).digest()
        claims = json.loads(decode(payload))

        if json.loads(decode(header)).get("alg") != "HS256" or not hmac.compare_digest(expected, decode(signature)):
            raise ValueError("signature")

        if claims["exp"] <= time.time() or claims.get("iss") != "robust" or "credit" not in claims.get("aud", []) or claims.get("role") != "SYSTEM":
            raise ValueError("claims")
    except (ValueError, KeyError, TypeError):
        context.abort(grpc.StatusCode.UNAUTHENTICATED, "Invalid service identity")


def affordability_message(a):
    if a is None:
        return None
    return pb.Affordability(status=a.status, cfads_minor=a.cfads_minor, max_payment_minor=a.max_payment_minor, max_amount_minor=a.max_amount_minor, dscr=a.dscr, reason=a.reason)


class Scoring(rpc.ScoringServiceServicer):
    def __init__(self):
        self.model = RiskModel()

    def Evaluate(self, request, context):
        authenticate(context)
        f = request.features

        if not request.scoring_request_id or not request.application_id or request.application_version < 1 or f.term_months < 1 or f.term_months > 180 or f.requested_amount_minor <= 0:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "Invalid scoring request")

        if f.business_age_months < 0 or any(value < 0 or value > 10000000000000 for value in (f.monthly_revenue_minor, f.monthly_expenses_minor, f.monthly_debt_payments_minor)):
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "Invalid features")

        risk = self.model.evaluate(f, request.product if request.HasField("product") else None)

        return pb.ScoringOutput(scoring_request_id=request.scoring_request_id, application_id=request.application_id, application_version=request.application_version, scoring=pb.Scoring(score=risk.score, risk_level=risk.level, factors=risk.factors, model_version=self.model.version, recommendation=risk.recommendation, affordability=affordability_message(risk.affordability)))


def main():
    address = os.getenv("GRPC_ADDR", "[::]:9093")

    if len(sys.argv) > 1 and sys.argv[1] == "healthcheck":
        with grpc.insecure_channel("127.0.0.1:" + address.rsplit(":", 1)[1]) as channel:
            reply = health_pb2_grpc.HealthStub(channel).Check(health_pb2.HealthCheckRequest(), timeout=2)

            if reply.status != health_pb2.HealthCheckResponse.SERVING:
                raise SystemExit(1)

        return

    if len(os.environ.get("JWT_SECRET", "")) < 32:
        raise RuntimeError("JWT_SECRET must contain at least 32 bytes")

    logging.basicConfig(level=logging.INFO, format="%(message)s")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=8), options=(("grpc.max_receive_message_length", 1048576),))
    rpc.add_ScoringServiceServicer_to_server(Scoring(), server)
    health_service = health.HealthServicer()
    health_service.set("", health_pb2.HealthCheckResponse.SERVING)
    health_pb2_grpc.add_HealthServicer_to_server(health_service, server)

    if not server.add_insecure_port(address):
        raise RuntimeError("Cannot bind scoring port")

    server.start()
    signal.signal(signal.SIGTERM, lambda *_: server.stop(10))
    signal.signal(signal.SIGINT, lambda *_: server.stop(10))
    logging.info("scoring listening on %s", address)
    server.wait_for_termination()

if __name__ == "__main__":
    main()
