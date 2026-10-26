import hashlib
import json
import math
import random
from datetime import date, timedelta
from pathlib import Path

REFERENCE_DATE = date(2026, 10, 1)

OKVED_SECTIONS = {"G": 0.35, "F": 0.12, "C": 0.10, "H": 0.10, "I": 0.08, "M": 0.10, "J": 0.07, "S": 0.08}

FIXED = json.loads((Path(__file__).parent / "profiles.json").read_text(encoding="utf-8"))


class InvalidINN(ValueError):
    pass


def legal_form(inn: str) -> str:
    if not inn.isdigit() or len(inn) not in (10, 12):
        raise InvalidINN(inn)
    return "OOO" if len(inn) == 10 else "IP"


def _rng(inn: str, source: str) -> random.Random:
    digest = hashlib.sha256(f"{source}:{inn}".encode()).hexdigest()
    return random.Random(int(digest[:16], 16))


def _quality(inn: str) -> float:
    return _rng(inn, "quality").gauss(0, 1)


def _sigmoid(x: float) -> float:
    return 1 / (1 + math.exp(-x))


def _poisson(rng: random.Random, rate: float) -> int:
    limit, k, p = math.exp(-rate), 0, 1.0
    while True:
        p *= rng.random()
        if p <= limit:
            return k
        k += 1


def _history(rng: random.Random, q: float, base_rate: float, q_effect: float) -> dict:
    overdues = min(_poisson(rng, math.exp(base_rate + q_effect * q)), 10)
    if overdues:
        max_days = min(365, round(30 + rng.expovariate(1 / (25 * math.exp(0.6 * q)))))
    else:
        max_days = rng.randrange(0, 30)
    return {"found": True, "overdueCount12m": overdues, "maxOverdueDays24m": max_days,
            "activeLoans": _poisson(rng, 1.2)}


def _not_found() -> dict:
    return {"found": False, "overdueCount12m": 0, "maxOverdueDays24m": 0, "activeLoans": 0}


def registry(inn: str) -> dict:
    form = legal_form(inn)
    fixed = FIXED.get(inn, {}).get("fns")
    if fixed is not None:
        return {"inn": inn, "legalForm": form, **fixed}

    rng, q = _rng(inn, "fns"), _quality(inn)
    age_months = max(1, min(300, round(rng.expovariate(1 / 50))))
    has_tax_debt = rng.random() < _sigmoid(-2.8 + 0.8 * q + 0.5 * (age_months < 12))
    return {
        "inn": inn,
        "status": "ACTIVE",
        "legalForm": form,
        "okvedSection": rng.choices(list(OKVED_SECTIONS), weights=list(OKVED_SECTIONS.values()))[0],
        "registeredAt": (REFERENCE_DATE - timedelta(days=round(age_months * 30.4))).isoformat(),
        "taxDebtMinor": round(rng.lognormvariate(math.log(50_000), 1.0)) * 100 if has_tax_debt else 0,
    }


def credit_report(inn: str) -> dict:
    form = legal_form(inn)
    fixed = FIXED.get(inn, {}).get("bki")
    if fixed is not None:
        return {"inn": inn, **fixed}

    rng, q = _rng(inn, "bki"), _quality(inn)
    company = None
    if form == "OOO":
        company = _history(rng, q, -2.0, 0.9) if rng.random() < 0.75 else _not_found()
    owner = _history(rng, q, -2.2, 0.8) if rng.random() < 0.85 else _not_found()
    return {"inn": inn, "company": company, "owner": owner}


def account(inn: str) -> dict:
    legal_form(inn)
    fixed = FIXED.get(inn, {}).get("bank")
    if fixed is not None:
        return {"inn": inn, **fixed}

    rng = _rng(inn, "bank")
    if rng.random() >= 0.4:
        return {"inn": inn, "found": False, "turnoverMonthlyMinor": 0}
    return {"inn": inn, "found": True, "turnoverMonthlyMinor": round(rng.lognormvariate(math.log(1_500_000), 0.8)) * 100}


def fault(inn: str, source: str) -> str | None:
    return FIXED.get(inn, {}).get("faults", {}).get(source)
