import json
import os
import time
import uuid
from pathlib import Path
from urllib.error import HTTPError
from urllib.request import Request, urlopen

BASE = os.getenv("API_URL", "http://localhost:8080/api/v1")
config = {}
config_file = Path(__file__).resolve().parents[1] / ".env"

if config_file.exists():
    config = dict(line.split("=", 1) for line in config_file.read_text().splitlines() if "=" in line and not line.lstrip().startswith(chr(35)))

PASSWORD = os.getenv("DEMO_PASSWORD", config.get("DEMO_PASSWORD", "Demo-Credit-2026!"))
PDF = b"%PDF-1.4\n1 0 obj << /Type /Catalog >> endobj\ntrailer << /Root 1 0 R >>\n%%EOF\n"


def request(method, path, token=None, body=None, version=None, key=None, expected=200, raw=None, content_type=None):
    headers = {}

    if token:
        headers["Authorization"] = "Bearer " + token

    if method in ("POST", "PATCH", "DELETE"):
        headers["Idempotency-Key"] = key or str(uuid.uuid4())

    if version is not None:
        headers["If-Match"] = '"' + str(version) + '"'

    data = raw

    if body is not None:
        data = json.dumps(body).encode()
        headers["Content-Type"] = "application/json"

    if content_type:
        headers["Content-Type"] = content_type

    req = Request(BASE + path, data=data, headers=headers, method=method)

    try:
        response = urlopen(req, timeout=30)
    except HTTPError as error:
        response = error

    status = response.code
    payload = response.read()
    assert status == expected, (method, path, status, expected, payload.decode(errors="replace"))

    if not payload:
        return None

    if "application/json" in response.headers.get("Content-Type", ""):
        body = json.loads(payload)

        if os.getenv("CHECK_CONTRACT") == "1":
            from validate_contract import check_response
            check_response(method, path, status, body)

        return body

    return payload


def login(name):
    return request("POST", "/auth/login", body={"login": name + "@demo.local", "password": PASSWORD})["accessToken"]


def company(token, age=36, revenue=300000000, overdue=False):
    payload = {"inn": "7700000000", "name": "Синтетическая компания " + uuid.uuid4().hex[:6], "businessAgeMonths": age, "monthlyRevenueMinor": revenue, "monthlyExpensesMinor": revenue // 2, "monthlyDebtPaymentsMinor": 0, "hasOverdueDebt": overdue}
    c = request("POST", "/companies", token, payload, expected=201)
    request("POST", f'/companies/{c["id"]}/verification', token)

    return c


def fresh(token, app):
    return request("GET", f'/applications/{app["id"]}', token)


def command(token, app, action, body=None, expected=200, key=None, version=None):
    current = fresh(token, app)

    return request("POST", f'/applications/{app["id"]}/{action}', token, body, current["version"] if version is None else version, key, expected)


def upload(token, app, document_type, item_id="", content=PDF, mime="application/pdf", expected=201, key=None):
    boundary = "boundary" + uuid.uuid4().hex
    chunks = []

    for field, value in (("type", document_type), ("requestItemId", item_id)):
        chunks.append(f'--{boundary}\r\nContent-Disposition: form-data; name="{field}"\r\n\r\n{value}\r\n'.encode())

    chunks.append(f'--{boundary}\r\nContent-Disposition: form-data; name="file"; filename="synthetic.pdf"\r\nContent-Type: {mime}\r\n\r\n'.encode() + content + b"\r\n")
    chunks.append(f"--{boundary}--\r\n".encode())

    return request("POST", f'/applications/{app["id"]}/documents', token, version=fresh(token, app)["version"], key=key, expected=expected, raw=b"".join(chunks), content_type="multipart/form-data; boundary=" + boundary)


def draft(token, c, amount=100000000, docs=True):
    a = request("POST", "/applications", token, {"need": {"companyId": c["id"], "amountMinor": amount, "termMonths": 12, "purpose": "WORKING_CAPITAL"}, "productId": "unsecured", "consents": {"personalDataProcessing": True, "creditHistoryCheck": True}}, expected=201)

    if docs:
        for typ in ("FINANCIAL_REPORT", "BANK_STATEMENT"):
            upload(token, a, typ)

    return a


def wait(token, app, statuses):
    until = time.monotonic() + 35

    while time.monotonic() < until:
        a = fresh(token, app)

        if a["status"] in statuses:
            return a

        time.sleep(0.5)

    raise AssertionError(("processing timeout", a))


def run():
    client, employee, other = login("client"), login("employee"), login("other")
    new_login = "e2e-signup-" + uuid.uuid4().hex + "@demo.local"
    credentials = {"login": new_login, "displayName": "Клиент E2E", "password": "E2E-Register-2026!"}
    session = request("POST", "/auth/register", body=credentials, expected=201)
    assert session["role"] == "CLIENT"
    request("POST", "/auth/register", body=credentials, expected=409)
    signed_in = request("POST", "/auth/login", body={"login": new_login.upper(), "password": credentials["password"]})
    client = signed_in["accessToken"]
    assert request("GET", "/me", client)["id"] == request("GET", "/me", session["accessToken"])["id"]
    print("PASS registration, duplicate prevention and login")
    request("GET", "/products", expected=401)
    request("GET", "/bank/applications", client, expected=403)
    c = company(client)
    request("GET", f'/companies/{c["id"]}', other, expected=404)
    recommendations = request("POST", "/recommendations", client, {"companyId": c["id"], "termMonths": 12, "purpose": "WORKING_CAPITAL", "collateral": "NONE"})
    assert any(x["availability"] == "BEST_MATCH" for x in recommendations)
    calc = request("POST", "/calculations", client, {"productId": "unsecured", "amountMinor": 100000000, "termMonths": 12})
    assert sum(x["principalMinor"] for x in calc["schedule"]) == 100000000
    assert calc["schedule"][-1]["balanceMinor"] == 0
    a = draft(client, c, docs=False)
    request("GET", f'/applications/{a["id"]}', other, expected=404)
    command(client, a, "submit", expected=422)
    upload(client, a, "OTHER", content=b"not a pdf", expected=415)
    file_key = str(uuid.uuid4())
    doc = upload(client, a, "FINANCIAL_REPORT", key=file_key)
    downloaded = request("GET", f'/applications/{a["id"]}/documents/{doc["id"]}/download', client)
    assert downloaded == PDF
    request("GET", f'/applications/{a["id"]}/documents/{doc["id"]}/download', other, expected=404)
    upload(client, a, "BANK_STATEMENT")
    extra = upload(client, a, "OTHER")
    before_delete = fresh(client, a)
    delete_key = str(uuid.uuid4())
    request("DELETE", f'/applications/{a["id"]}/documents/{extra["id"]}', client, version=before_delete["version"], key=delete_key, expected=204)
    request("DELETE", f'/applications/{a["id"]}/documents/{extra["id"]}', client, version=before_delete["version"], key=delete_key, expected=204)
    request("GET", f'/applications/{a["id"]}/documents/{extra["id"]}/download', client, expected=404)
    command(client, a, "validation")
    command(client, a, "submit", expected=202)
    a = wait(client, a, {"APPROVED"})
    decision = request("GET", f'/applications/{a["id"]}/decision', client)
    command(client, a, "accept", {"offerId": decision["offers"][0]["id"]})
    accepted = fresh(client, a)
    key = str(uuid.uuid4())
    first = command(client, a, "issuance", expected=201, key=key, version=accepted["version"])
    second = command(client, a, "issuance", expected=201, key=key, version=accepted["version"])
    assert first["id"] == second["id"] and first["simulated"]
    command(client, a, "issuance", expected=409)
    assert fresh(client, a)["status"] == "ISSUED"
    replay = upload(client, a, "FINANCIAL_REPORT", key=file_key)
    assert replay["id"] == doc["id"]
    upload(client, a, "FINANCIAL_REPORT", key=file_key, content=PDF + b"different", expected=409)
    assert fresh(client, a)["status"] == "ISSUED"
    print("PASS approval, file roundtrip, ownership, idempotent issuance and upload replay")
    alternative_company = company(client, revenue=30000000)
    alt = draft(client, alternative_company, amount=500000000)
    command(client, alt, "submit", expected=202)
    wait(client, alt, {"ALTERNATIVE"})
    d = request("GET", f'/applications/{alt["id"]}/decision', client)
    assert d["offers"][0]["amountMinor"] < 500000000
    command(client, alt, "decline", {"reason": "Не подходит"})
    print("PASS alternative offer and client decline")
    manual_company = company(client, age=6, overdue=True)
    manual = draft(client, manual_company)
    command(client, manual, "submit", expected=202)
    wait(client, manual, {"MANUAL_CHECK"})
    bank = request("GET", f'/bank/applications/{manual["id"]}', employee)
    assert bank["scoring"]["modelVersion"] == "rules-v1"
    current = fresh(client, manual)
    request("POST", f'/bank/applications/{manual["id"]}/decision', employee, {"outcome": "NEEDS_INFO", "reason": "Дополните сведения", "items": [{"title": "Дополнительная выписка", "type": "OTHER", "required": True}]}, current["version"])
    info = request("GET", f'/applications/{manual["id"]}/information-requests', client)
    command(client, manual, "submit", expected=422)
    upload(client, manual, "OTHER", info[0]["items"][0]["id"])
    command(client, manual, "submit", expected=202)
    wait(client, manual, {"MANUAL_CHECK"})
    current = fresh(client, manual)
    request("POST", f'/bank/applications/{manual["id"]}/decision', employee, {"outcome": "REJECTED", "reason": "Высокая долговая нагрузка", "offers": []}, current["version"])
    assert fresh(client, manual)["status"] == "REJECTED"
    command(client, manual, "accept", {"offerId": "not-an-offer"}, expected=409)
    manual_offer = draft(client, manual_company)
    command(client, manual_offer, "submit", expected=202)
    wait(client, manual_offer, {"MANUAL_CHECK"})
    current = fresh(client, manual_offer)
    request("POST", f'/bank/applications/{manual_offer["id"]}/decision', employee, {"outcome": "ALTERNATIVE", "reason": "Снижена сумма", "offers": [{"amountMinor": 80000000, "termMonths": 12, "annualRateBps": 1800}]}, current["version"])
    offered = request("GET", f'/applications/{manual_offer["id"]}/decision', client)
    assert offered["offers"][0]["expiresAt"] and offered["offers"][0]["monthlyPaymentMinor"] > 0
    command(client, manual_offer, "accept", {"offerId": offered["offers"][0]["id"]})
    command(client, manual_offer, "issuance", expected=201)
    print("PASS manual review, information request, resubmission, rejection and manual offer")
    cancel = draft(client, c)
    command(client, cancel, "submit", expected=202)
    command(client, cancel, "cancel")
    time.sleep(2)
    assert fresh(client, cancel)["status"] == "CANCELLED"
    request("PATCH", f'/applications/{cancel["id"]}', client, {}, version=1, expected=409)
    feed = request("GET", "/feed", client)
    assert any(x.get("applicationId") == a["id"] for x in feed)
    print("PASS cancellation and personalized feed")
    print("ALL E2E CHECKS PASSED")

if __name__ == "__main__":
    run()
