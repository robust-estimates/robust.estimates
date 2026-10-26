import argparse
import hashlib
import json
import os
import re
import sys
import time
import uuid
from pathlib import Path
from urllib.error import HTTPError, URLError
from urllib.parse import urlencode, urlsplit
from urllib.request import Request, urlopen

PDF = b"%PDF-1.4\n1 0 obj << /Type /Catalog >> endobj\ntrailer << /Root 1 0 R >>\n%%EOF\n"
BASE = ""
PASSWORD = ""
EVENTS = []
TOKENS = []
VERBOSE = False
WAIT_SECONDS = 45
EXPECTED_ROUTES = [('POST', '/auth/login'), ('POST', '/auth/register'), ('GET', '/me'), ('POST', '/companies'), ('GET', '/companies'), ('GET', '/companies/{companyId}'), ('PATCH', '/companies/{companyId}'), ('POST', '/companies/{companyId}/verification'), ('GET', '/products'), ('GET', '/products/{productId}'), ('POST', '/recommendations'), ('POST', '/calculations'), ('POST', '/applications'), ('GET', '/applications'), ('GET', '/applications/{applicationId}'), ('PATCH', '/applications/{applicationId}'), ('POST', '/applications/{applicationId}/validation'), ('POST', '/applications/{applicationId}/submit'), ('POST', '/applications/{applicationId}/cancel'), ('GET', '/applications/{applicationId}/decision'), ('POST', '/applications/{applicationId}/accept'), ('POST', '/applications/{applicationId}/decline'), ('POST', '/applications/{applicationId}/issuance'), ('GET', '/applications/{applicationId}/history'), ('GET', '/applications/{applicationId}/information-requests'), ('GET', '/applications/{applicationId}/documents'), ('POST', '/applications/{applicationId}/documents'), ('DELETE', '/applications/{applicationId}/documents/{documentId}'), ('GET', '/applications/{applicationId}/documents/{documentId}/download'), ('GET', '/bank/applications'), ('GET', '/bank/applications/{applicationId}/scoring'), ('GET', '/bank/applications/{applicationId}'), ('POST', '/bank/applications/{applicationId}/decision'), ('GET', '/feed')]


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

        return body

    return payload


def login(name):
    return request("POST", "/auth/login", body={"login": name + "@demo.local", "password": PASSWORD})["accessToken"]


def company(token, age=36, revenue=300000000, inn="7700000000"):
    payload = {"inn": inn, "name": "Синтетическая компания " + uuid.uuid4().hex[:6], "businessAgeMonths": age, "monthlyRevenueMinor": revenue, "monthlyExpensesMinor": revenue // 2, "monthlyDebtPaymentsMinor": 0}
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
    until = time.monotonic() + WAIT_SECONDS

    while time.monotonic() < until:
        a = fresh(token, app)

        if a["status"] in statuses:
            return a

        time.sleep(0.5)

    raise AssertionError(("processing timeout", a))


raw_request = request


def request(method, path, token=None, body=None, version=None, key=None, expected=200, raw=None, content_type=None):
    started = time.monotonic()
    event = {"method": method, "path": path, "expected": expected}

    try:
        result = raw_request(method, path, token, body, version, key, expected, raw, content_type)
        event.update(status=expected, ok=True)

        if VERBOSE:
            print(f"    {method} {path} → {expected}", flush=True)

        return result
    except Exception as error:
        event.update(ok=False, error=safe(str(error)))

        if isinstance(error, AssertionError) and error.args and isinstance(error.args[0], tuple) and len(error.args[0]) >= 4:
            event["status"] = error.args[0][2]

        raise
    finally:
        event["durationMs"] = round((time.monotonic() - started) * 1000)
        EVENTS.append(event)


def safe(text):
    for secret in [PASSWORD, *TOKENS]:
        if secret:
            text = text.replace(secret, "[скрыто]")

    return text[:1800]


def check(condition, message):
    if not condition:
        raise AssertionError(message)


def load_password(path):
    if os.getenv("DEMO_PASSWORD"):
        return os.environ["DEMO_PASSWORD"]

    if path is not None:
        candidates = [Path(path)]

        if not candidates[0].is_file():
            raise ValueError("Файл окружения не найден: " + str(path))
    else:
        candidates = [Path.cwd() / ".env", Path(__file__).resolve().parents[1] / ".env"]

    for candidate in candidates:
        if candidate.is_file():
            for line in candidate.read_text().splitlines():
                line = line.strip()

                if line.startswith("export "):
                    line = line[7:].lstrip()

                if "=" not in line or line.startswith(chr(35)):
                    continue

                key, value = line.split("=", 1)

                if key.strip() != "DEMO_PASSWORD":
                    continue

                value = value.strip()

                if len(value) >= 2 and value[0] == value[-1] and value[0] in (chr(34), chr(39)):
                    return value[1:-1]

                return value.split(" #", 1)[0].rstrip()

    return "Demo-Credit-2026!"


def availability():
    root = BASE[:-7] if BASE.endswith("/api/v1") else None

    if root is not None:
        for path in ("/healthz", "/readyz"):
            with urlopen(Request(root + path), timeout=5) as response:
                check(response.status == 200, "Сервис не готов: " + path)
    else:
        request("GET", "/products", expected=401)


def company_payload(name, age=36, revenue=300000000, inn="7700000000"):
    return {"inn": inn, "name": name + " " + uuid.uuid4().hex[:6], "businessAgeMonths": age, "monthlyRevenueMinor": revenue, "monthlyExpensesMinor": revenue // 2, "monthlyDebtPaymentsMinor": 0}


def auth_test(state, args):
    for name, role in (("client", "CLIENT"), ("employee", "EMPLOYEE"), ("other", "CLIENT")):
        token = login(name)
        TOKENS.append(token)
        state[name] = token
        me = request("GET", "/me", token)
        check(me["role"] == role, "Неверная роль для " + name)

    request("GET", "/products", expected=401)
    request("GET", "/products", token="invalid-token", expected=401)
    request("GET", "/bank/applications", state["client"], expected=403)


def registration_test(state, args):
    login_value = "console-signup-" + uuid.uuid4().hex + "@demo.local"
    password = "Console-Register-2026!"
    body = {"login": login_value, "displayName": "Новый клиент", "password": password}
    session = request("POST", "/auth/register", body=body, expected=201)
    check(session["role"] == "CLIENT" and session["accessToken"], "Регистрация не вернула клиентскую сессию")
    TOKENS.append(session["accessToken"])
    me = request("GET", "/me", session["accessToken"])
    check(me["role"] == "CLIENT" and me["displayName"] == body["displayName"] and not me["companyIds"], "Профиль нового клиента некорректен")
    logged_in = request("POST", "/auth/login", body={"login": " " + login_value.upper() + " ", "password": password})
    TOKENS.append(logged_in["accessToken"])
    same = request("GET", "/me", logged_in["accessToken"])
    check(me["id"] == same["id"], "Повторный вход открыл другую учётную запись")
    request("GET", "/bank/applications", session["accessToken"], expected=403)
    request("POST", "/auth/register", body={**body, "login": login_value.upper()}, expected=409)
    request("POST", "/auth/register", body={**body, "login": "new-" + login_value, "role": "EMPLOYEE"}, expected=400)
    request("POST", "/auth/register", body={**body, "login": "invalid-email"}, expected=422)
    request("POST", "/auth/register", body={**body, "password": "short"}, expected=422)


def companies_test(state, args):
    client, other, employee = state["client"], state["other"], state["employee"]
    payload = company_payload("Console check")
    key = str(uuid.uuid4())
    c = request("POST", "/companies", client, payload, key=key, expected=201)
    repeated = request("POST", "/companies", client, payload, key=key, expected=201)
    check(c["id"] == repeated["id"], "Повтор создания компании вернул другой ID")
    request("POST", "/companies", client, {**payload, "name": payload["name"] + " changed"}, key=key, expected=409)
    request("GET", f'/companies/{c["id"]}', other, expected=404)
    request("POST", "/companies", employee, payload, expected=403)
    companies = request("GET", "/companies", client)
    check(any(item["id"] == c["id"] for item in companies), "Компания отсутствует в списке")
    payload["representativeName"] = "Тестовый представитель"
    request("PATCH", f'/companies/{c["id"]}', client, payload)
    request("POST", f'/companies/{c["id"]}/verification', client)
    c = request("GET", f'/companies/{c["id"]}', client)
    check(c["verificationStatus"] == "VERIFIED", "Компания не проверена")
    check(c["representativeName"] == payload["representativeName"], "Обновление профиля не сохранилось")
    state["company"] = c


def products_test(state, args):
    client, c = state["client"], state["company"]
    products = request("GET", "/products", client)
    check(any(p["id"] == "unsecured" for p in products), "Нет тестового продукта unsecured")
    request("GET", "/products/unsecured", client)
    request("GET", "/products/not-found-console-test", client, expected=404)
    request("GET", "/products?" + urlencode({"search": "без залога"}), client)
    recommendations = request("POST", "/recommendations", client, {"companyId": c["id"], "termMonths": 12, "purpose": "WORKING_CAPITAL", "collateral": "NONE"})
    check(any(p["availability"] == "BEST_MATCH" for p in recommendations), "Нет подходящих рекомендаций")
    calc = request("POST", "/calculations", client, {"productId": "unsecured", "amountMinor": 100000000, "termMonths": 12})
    check(sum(p["principalMinor"] for p in calc["schedule"]) == 100000000, "График не возвращает всю основную сумму")
    check(calc["schedule"][-1]["balanceMinor"] == 0, "В конце графика остался долг")
    request("POST", "/calculations", client, {"productId": "unsecured", "amountMinor": 100000000, "termMonths": 0}, expected=422)


def draft_test(state, args):
    client, other = state["client"], state["other"]
    a = draft(client, state["company"], docs=False)
    state["application"] = a
    request("GET", f'/applications/{a["id"]}', other, expected=404)
    body = {"need": {**a["need"], "description": "Проверка сохранения черновика"}}
    updated = request("PATCH", f'/applications/{a["id"]}', client, body, version=a["version"])
    check(updated["need"]["description"] == body["need"]["description"], "Изменение черновика не сохранилось")
    request("PATCH", f'/applications/{a["id"]}', client, body, version=a["version"], expected=409)
    page = request("GET", "/applications?" + urlencode({"companyId": state["company"]["id"], "limit": 100}), client)
    check(any(row["id"] == a["id"] for row in page["items"]), "Заявка отсутствует в списке")
    validation = command(client, a, "validation")
    check(not validation["valid"], "Заявка без документов прошла проверку")
    command(client, a, "submit", expected=422)
    request("GET", f'/applications/{a["id"]}/decision', client, expected=409)


def files_test(state, args):
    client, other, a = state["client"], state["other"], state["application"]
    upload(client, a, "OTHER", content=b"not a pdf", expected=415)
    file_key = str(uuid.uuid4())
    doc = upload(client, a, "FINANCIAL_REPORT", key=file_key)
    state["file_key"], state["document"] = file_key, doc
    upload(client, a, "BANK_STATEMENT")
    docs = request("GET", f'/applications/{a["id"]}/documents', client)
    check(len(docs) == 2, "Неверное число документов")
    content = request("GET", f'/applications/{a["id"]}/documents/{doc["id"]}/download', client)
    check(content == PDF, "Скачанный файл отличается от загруженного")
    check(doc["checksum"] == hashlib.sha256(PDF).hexdigest(), "SHA-256 документа не совпадает")
    request("GET", f'/applications/{a["id"]}/documents/{doc["id"]}/download', other, expected=404)
    extra = upload(client, a, "OTHER")
    version = fresh(client, a)["version"]
    key = str(uuid.uuid4())

    for _ in range(2):
        request("DELETE", f'/applications/{a["id"]}/documents/{extra["id"]}', client, version=version, key=key, expected=204)

    request("GET", f'/applications/{a["id"]}/documents/{extra["id"]}/download', client, expected=404)

    if args.extended:
        upload(client, a, "OTHER", content=PDF + b"0" * (10 * 1024 * 1024), expected=413)

    check(command(client, a, "validation")["valid"], "Полный комплект документов не прошёл проверку")


def approval_test(state, args):
    client, employee, other, a = state["client"], state["employee"], state["other"], state["application"]
    command(client, a, "submit", expected=202)
    wait(client, a, {"APPROVED"})
    scoring = request("GET", f'/bank/applications/{a["id"]}/scoring', employee)
    check(scoring["riskLevel"] == "LOW", "Для контрольной компании ожидался низкий риск")
    request("GET", f'/bank/applications/{a["id"]}/scoring', client, expected=403)
    request("GET", f'/bank/applications/{a["id"]}', employee)
    queue = request("GET", "/bank/applications?status=APPROVED&limit=100&" + urlencode({"search": a["number"]}), employee)
    check(any(row["id"] == a["id"] for row in queue["items"]), "Заявки нет в банковской очереди")
    decision = request("GET", f'/applications/{a["id"]}/decision', client)
    command(client, a, "accept", {"offerId": decision["offers"][0]["id"]})
    version = fresh(client, a)["version"]
    key = str(uuid.uuid4())
    first = command(client, a, "issuance", expected=201, key=key, version=version)
    repeat = command(client, a, "issuance", expected=201, key=key, version=version)
    check(first["id"] == repeat["id"] and first["simulated"], "Повтор выдачи создал другой кредит")
    command(client, a, "issuance", expected=409)
    replay = upload(client, a, "FINANCIAL_REPORT", key=state["file_key"])
    check(replay["id"] == state["document"]["id"], "Повтор загрузки создал другой документ")
    upload(client, a, "FINANCIAL_REPORT", key=state["file_key"], content=PDF + b"changed", expected=409)
    check(fresh(client, a)["status"] == "ISSUED", "Итоговый статус не ISSUED")
    history = request("GET", f'/applications/{a["id"]}/history', client)
    statuses = [e["toStatus"] for e in history]

    for value in ("IN_PROGRESS", "SCORING", "APPROVED", "ACCEPTED", "ISSUED"):
        check(value in statuses, "В истории отсутствует статус " + value)

    feed = request("GET", "/feed", client)
    check(any(x.get("applicationId") == a["id"] for x in feed), "В ленте нет результата заявки")
    other_feed = request("GET", "/feed", other)
    check(not any(x.get("applicationId") == a["id"] for x in other_feed), "Чужая заявка попала в ленту")


def alternative_test(state, args):
    client = state["client"]
    c = company(client, revenue=30000000)
    a = draft(client, c, amount=500000000)
    command(client, a, "submit", expected=202)
    wait(client, a, {"ALTERNATIVE"})
    decision = request("GET", f'/applications/{a["id"]}/decision', client)
    check(decision["offers"][0]["amountMinor"] < 500000000, "Альтернативная сумма не уменьшена")
    command(client, a, "decline", {"reason": "Проверка отказа клиента"})
    check(fresh(client, a)["status"] == "DECLINED", "Отказ клиента не сохранён")


def manual_test(state, args):
    client, employee = state["client"], state["employee"]
    c = company(client, age=6, inn="0000090029")
    state["manual_company"] = c
    a = draft(client, c)
    command(client, a, "submit", expected=202)
    wait(client, a, {"MANUAL_CHECK"})
    manual_body = {"outcome": "NEEDS_INFO", "reason": "Дополните сведения", "items": [{"title": "Дополнительная выписка", "type": "OTHER", "required": True}]}
    current = fresh(client, a)
    path = f'/bank/applications/{a["id"]}/decision'
    request("POST", path, client, manual_body, version=current["version"], expected=403)
    request("POST", path, employee, manual_body, version=current["version"])
    info = request("GET", f'/applications/{a["id"]}/information-requests', client)
    check(info and info[0]["items"], "Запрос сведений не сохранён")
    command(client, a, "submit", expected=422)
    upload(client, a, "OTHER", info[0]["items"][0]["id"])
    command(client, a, "submit", expected=202)
    wait(client, a, {"MANUAL_CHECK"})
    request("POST", path, employee, {"outcome": "REJECTED", "reason": "Демонстрационный отказ", "offers": []}, version=fresh(client, a)["version"])
    check(fresh(client, a)["status"] == "REJECTED", "Отказ банка не сохранён")
    command(client, a, "accept", {"offerId": "missing-offer"}, expected=409)


def manual_offer_test(state, args):
    client, employee = state["client"], state["employee"]

    for outcome, amount in (("ALTERNATIVE", 80000000), ("APPROVED", 100000000)):
        a = draft(client, state["manual_company"])
        command(client, a, "submit", expected=202)
        wait(client, a, {"MANUAL_CHECK"})
        body = {"outcome": outcome, "reason": "Решение специалиста в тесте", "offers": [{"amountMinor": amount, "termMonths": 12, "annualRateBps": 1800}]}
        request("POST", f'/bank/applications/{a["id"]}/decision', employee, body, version=fresh(client, a)["version"])
        decision = request("GET", f'/applications/{a["id"]}/decision', client)
        offer = decision["offers"][0]
        check(offer["expiresAt"] and offer["monthlyPaymentMinor"] > 0, "Сервер не рассчитал предложение")
        command(client, a, "accept", {"offerId": offer["id"]})
        command(client, a, "issuance", expected=201)
        check(fresh(client, a)["status"] == "ISSUED", "Ручное решение не дошло до выдачи")


def cancellation_test(state, args):
    client = state["client"]
    a = draft(client, state["manual_company"])
    command(client, a, "submit", expected=202)
    wait(client, a, {"MANUAL_CHECK"})
    command(client, a, "cancel")
    time.sleep(2)
    check(fresh(client, a)["status"] == "CANCELLED", "Фоновая обработка изменила отменённую заявку")
    command(client, a, "submit", expected=409)


def coverage():
    covered = set()

    for event in EVENTS:
        if not event.get("ok"):
            continue

        path = event["path"].split("?", 1)[0]

        for method, template in EXPECTED_ROUTES:
            if method == event["method"] and re.fullmatch(re.sub(r"\{[^}]+\}", "[^/]+", template), path):
                covered.add((method, template))

    return covered


def main():
    global BASE, PASSWORD, VERBOSE, WAIT_SECONDS
    parser = argparse.ArgumentParser(description="Консольная проверка API robust.estimates. Создаёт только синтетические данные, существующие данные не удаляет.")
    parser.add_argument("--base-url", default=os.getenv("API_URL", "http://localhost:8080/api/v1"))
    parser.add_argument("--env-file")
    parser.add_argument("--verbose", action="store_true", help="Показать каждый HTTP-запрос без токенов и паролей")
    parser.add_argument("--extended", action="store_true", help="Дополнительно проверить файл размером больше 10 МиБ")
    parser.add_argument("--wait-seconds", type=int, default=45, help="Максимум ожидания одного решения")
    parser.add_argument("--report", help="Сохранить JSON-отчёт без токенов и паролей")
    args = parser.parse_args()

    if not __debug__:
        parser.error("Не используйте Python с флагом -O")

    BASE = args.base_url.rstrip("/")
    parsed = urlsplit(BASE)

    if parsed.scheme not in ("http", "https") or not parsed.hostname or parsed.username or parsed.password or parsed.query or parsed.fragment:
        parser.error("Нужен HTTP(S)-адрес API без встроенных паролей и query-параметров")

    if args.wait_seconds < 1:
        parser.error("--wait-seconds должен быть больше нуля")

    try:
        PASSWORD = load_password(args.env_file)
    except (ValueError, OSError) as error:
        parser.error(str(error))

    VERBOSE, WAIT_SECONDS = args.verbose, args.wait_seconds
    tests = [("Готовность API", lambda state, args: availability()), ("Вход, роли, защита методов", auth_test), ("Регистрация клиента и запрет выбора роли", registration_test), ("Компании, профиль и повтор команд", companies_test), ("Каталог, подбор и калькулятор", products_test), ("Черновик, версии и комплектность", draft_test), ("Документы, S3 и доступ к файлам", files_test), ("Скоринг, одобрение, выдача и история", approval_test), ("Альтернативные условия и отказ клиента", alternative_test), ("Ручная проверка, запрос сведений и отказ банка", manual_test), ("Ручное одобрение и альтернативное предложение", manual_offer_test), ("Отмена и фоновая обработка", cancellation_test)]
    state, results = {}, []
    print("robust.estimates — консольная проверка", flush=True)
    print("API:", BASE, flush=True)
    print("Будут созданы синтетические компании, заявки и файлы. Старые данные сохраняются.\n", flush=True)
    started = time.monotonic()
    failed = False

    for number, (title, test) in enumerate(tests, 1):
        print(f"[{number}/{len(tests)}] {title}…", flush=True)
        step_started = time.monotonic()

        try:
            test(state, args)
            result = {"name": title, "ok": True, "durationSeconds": round(time.monotonic() - step_started, 2)}
            print("  OK", flush=True)
        except (Exception, KeyboardInterrupt) as error:
            failed = True
            message = "Проверка прервана" if isinstance(error, KeyboardInterrupt) else safe(str(error))

            if isinstance(error, (URLError, TimeoutError, ConnectionError)):
                message = "API недоступен или запрос превысил таймаут: " + message

            result = {"name": title, "ok": False, "error": message}
            print("  FAIL:", message, flush=True)

            if EVENTS:
                last = EVENTS[-1]
                print("  Последний запрос:", last["method"], last["path"], flush=True)

            print("  Дальнейшие шаги остановлены: они зависят от этого результата.", flush=True)

        results.append(result)

        if failed:
            break

    covered = coverage()
    missing = [f"{method} {path}" for method, path in EXPECTED_ROUTES if (method, path) not in covered]

    if not failed and missing:
        failed = True
        print("Не проверены операции:", ", ".join(missing), flush=True)

    elapsed = round(time.monotonic() - started, 2)
    print(f"\nИтог: {sum(item['ok'] for item in results)}/{len(tests)} этапов; API {len(covered)}/{len(EXPECTED_ROUTES)} операций; {len(EVENTS)} запросов; {elapsed} с.", flush=True)
    print("ПРОВЕРКА НЕ ПРОЙДЕНА" if failed else "ВСЕ ПРОВЕРКИ ПРОЙДЕНЫ", flush=True)
    report = {"baseUrl": BASE, "ok": not failed, "durationSeconds": elapsed, "steps": results, "coverage": {"checked": len(covered), "total": len(EXPECTED_ROUTES), "missing": missing}, "requests": EVENTS}

    if args.report:
        try:
            path = Path(args.report)
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n")
            print("Отчёт:", path, flush=True)
        except OSError as error:
            print("Не удалось сохранить отчёт:", safe(str(error)), file=sys.stderr)

            return 1

    return 1 if failed else 0

if __name__ == "__main__":
    raise SystemExit(main())
