import argparse
from pathlib import Path

import numpy as np
import pandas as pd

from schema import Application, CreditHistory, LegalForm, Purpose

TARGET_DEFAULT_RATE = 0.06


OKVED_SECTIONS = {
    "G": (0.35, 0.0),
    "F": (0.12, 0.3),
    "C": (0.10, 0.0),
    "H": (0.10, 0.2),
    "I": (0.08, 0.4),
    "M": (0.10, -0.2),
    "J": (0.07, -0.3),
    "S": (0.08, 0.1),
}
REGIONS = ["77", "78", "50", "16", "66", "54", "23", "02"]
TERMS_MONTHS = [6, 12, 18, 24, 36, 48, 60]
PURPOSES = [Purpose.WORKING_CAPITAL, Purpose.EQUIPMENT, Purpose.BUSINESS_EXPANSION]


MIN_AMOUNT_MINOR = 10_000_000
MAX_AMOUNT_MINOR = 5_000_000_000


TRUE_LOGIT_COEFS = {
    "log_business_age_years": -0.4,
    "debt_load": 2.5,
    "margin": -2.0,
    "company_overdues": 0.5,
    "owner_overdues": 0.6,
    "overdue_90_plus": 0.8,
    "no_credit_history": 0.3,
    "tax_debt": 0.7,
    "bank_client": -0.2,
    "hidden_quality": 0.4,
}


def _sigmoid(x: np.ndarray) -> np.ndarray:
    return 1 / (1 + np.exp(-x))


def _inn_check_digit(digits: str, weights: list[int]) -> str:
    return str(sum(int(d) * w for d, w in zip(digits, weights)) % 11 % 10)


def make_inn(index: int, legal_form: str) -> str:

    if legal_form == LegalForm.OOO:
        base = f"00{index:07d}"
        return base + _inn_check_digit(base, [2, 4, 10, 3, 5, 9, 4, 6, 8])
    base = f"00{index:08d}"
    base += _inn_check_digit(base, [7, 2, 4, 10, 3, 5, 9, 4, 6, 8])
    return base + _inn_check_digit(base, [3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8])


def _nullable_int(values: np.ndarray, available: np.ndarray) -> pd.api.extensions.ExtensionArray:

    result = pd.array(np.asarray(values).round().astype(np.int64), dtype="Int64")
    result[~available] = pd.NA
    return result


def _credit_history(rng, q, available, base_rate, q_effect, has_debt):
    n = len(q)
    overdues = np.clip(rng.poisson(np.exp(base_rate + q_effect * q)), 0, 10)
    long_overdue = 30 + rng.exponential(25 * np.exp(0.6 * q))
    max_days = np.where(overdues > 0, long_overdue, rng.integers(0, 30, n))
    max_days = np.clip(max_days, 0, 365).round()
    active = rng.poisson(1.2, n) + has_debt
    return (_nullable_int(overdues, available), _nullable_int(max_days, available),
            _nullable_int(active, available))


def _calibrate_intercept(logit_wo_intercept: np.ndarray, target_rate: float) -> float:
    low, high = -20.0, 20.0
    for _ in range(100):
        mid = (low + high) / 2
        if _sigmoid(mid + logit_wo_intercept).mean() < target_rate:
            low = mid
        else:
            high = mid
    return (low + high) / 2


def generate_companies(n: int, seed: int = 42, target_default_rate: float = TARGET_DEFAULT_RATE) -> pd.DataFrame:
    rng = np.random.default_rng(seed)
    q = rng.standard_normal(n)

    legal_form = np.where(rng.random(n) < 0.6, LegalForm.IP.value, LegalForm.OOO.value)
    sections = list(OKVED_SECTIONS)
    okved = rng.choice(sections, size=n, p=[OKVED_SECTIONS[s][0] for s in sections])
    region = rng.choice(REGIONS, size=n)
    age = np.clip(rng.exponential(50, n), 1, 300).round().astype(int)


    median_revenue_rub = np.where(legal_form == LegalForm.IP, 600_000, 2_500_000)
    log_revenue = np.log(median_revenue_rub) + 0.3 * (np.log1p(age / 12) - 1) + rng.normal(0, 0.8, n)
    revenue = (np.exp(log_revenue).round() * 100).astype(np.int64)

    margin = np.clip(rng.beta(4, 14, n) - 0.06 * q + rng.normal(0, 0.03, n), -0.3, 0.6)
    expenses = (revenue * (1 - margin)).round().astype(np.int64)

    has_debt = rng.random(n) < 0.55
    debt_load = np.where(has_debt, np.clip(rng.beta(2, 7, n) + 0.04 * q, 0.01, 0.9), 0.0)
    debt_payments = (revenue * debt_load).round().astype(np.int64)

    has_tax_debt = rng.random(n) < _sigmoid(-2.8 + 0.8 * q + 0.5 * (age < 12))
    tax_debt = np.where(has_tax_debt, rng.lognormal(np.log(50_000), 1.0, n).round() * 100, 0).astype(np.int64)



    company_available = (legal_form == LegalForm.OOO) & (age >= 6) & (rng.random(n) < 0.75)
    owner_available = rng.random(n) < 0.85
    company_overdues, company_max_days, company_active = _credit_history(
        rng, q, company_available, -2.0, 0.9, has_debt)
    owner_overdues, owner_max_days, owner_active = _credit_history(
        rng, q, owner_available, -2.2, 0.8, has_debt)

    bank_client = rng.random(n) < 0.4
    turnover = _nullable_int(revenue * rng.uniform(0.6, 1.0, n), bank_client)

    purpose = rng.choice([p.value for p in PURPOSES], size=n, p=[0.5, 0.3, 0.2])
    amount = np.clip(revenue * rng.uniform(1, 8, n), MIN_AMOUNT_MINOR, MAX_AMOUNT_MINOR)
    amount = (amount // 100_000 * 100_000).astype(np.int64)
    term = rng.choice(TERMS_MONTHS, size=n)
    has_collateral = rng.random(n) < 0.25

    c = TRUE_LOGIT_COEFS
    company_over = company_overdues.fillna(0).to_numpy(dtype=float)
    owner_over = owner_overdues.fillna(0).to_numpy(dtype=float)
    max_days_any = np.maximum(company_max_days.fillna(0).to_numpy(dtype=float),
                              owner_max_days.fillna(0).to_numpy(dtype=float))
    logit = (
        c["log_business_age_years"] * np.log1p(age / 12)
        + c["debt_load"] * debt_load
        + c["margin"] * margin
        + c["company_overdues"] * np.minimum(company_over, 4)
        + c["owner_overdues"] * np.minimum(owner_over, 4)
        + c["overdue_90_plus"] * (max_days_any >= 90)
        + c["no_credit_history"] * (~company_available & ~owner_available)
        + c["tax_debt"] * has_tax_debt
        + c["bank_client"] * bank_client
        + c["hidden_quality"] * q
        + np.array([OKVED_SECTIONS[s][1] for s in okved])
    )
    logit += _calibrate_intercept(logit, target_default_rate)
    true_pd = _sigmoid(logit)
    default = (rng.random(n) < true_pd).astype(int)

    return pd.DataFrame({
        "inn": [make_inn(i, f) for i, f in enumerate(legal_form, start=1)],
        "name": [f"ООО «Синтетика {i}»" if f == LegalForm.OOO else f"ИП Синтетический {i}"
                 for i, f in enumerate(legal_form, start=1)],
        "legal_form": legal_form,
        "okved_section": okved,
        "region_code": region,
        "business_age_months": age,
        "monthly_revenue_minor": revenue,
        "monthly_expenses_minor": expenses,
        "monthly_debt_payments_minor": debt_payments,
        "tax_debt_minor": tax_debt,
        "company_overdues_30_plus_12m": company_overdues,
        "company_max_overdue_days_24m": company_max_days,
        "company_active_loans": company_active,
        "owner_overdues_30_plus_12m": owner_overdues,
        "owner_max_overdue_days_24m": owner_max_days,
        "owner_active_loans": owner_active,
        "account_turnover_monthly_minor": turnover,
        "requested_amount_minor": amount,
        "requested_term_months": term,
        "purpose": purpose,
        "has_collateral": has_collateral,

        "true_pd": true_pd,
        "default_12m": default,
    })


def _history(row, prefix: str) -> CreditHistory | None:
    if pd.isna(row[f"{prefix}_overdues_30_plus_12m"]):
        return None
    return CreditHistory(
        overdues_30_plus_12m=int(row[f"{prefix}_overdues_30_plus_12m"]),
        max_overdue_days_24m=int(row[f"{prefix}_max_overdue_days_24m"]),
        active_loans=int(row[f"{prefix}_active_loans"]),
    )


def to_application(row) -> Application:

    turnover = row["account_turnover_monthly_minor"]
    return Application(
        business_age_months=int(row["business_age_months"]),
        monthly_revenue_minor=int(row["monthly_revenue_minor"]),
        monthly_expenses_minor=int(row["monthly_expenses_minor"]),
        monthly_debt_payments_minor=int(row["monthly_debt_payments_minor"]),
        requested_amount_minor=int(row["requested_amount_minor"]),
        requested_term_months=int(row["requested_term_months"]),
        purpose=Purpose(row["purpose"]),
        has_collateral=bool(row["has_collateral"]),
        legal_form=LegalForm(row["legal_form"]),
        okved_section=row["okved_section"],
        region_code=row["region_code"],
        tax_debt_minor=int(row["tax_debt_minor"]),
        company_credit_history=_history(row, "company"),
        owner_credit_history=_history(row, "owner"),
        credit_bureau_checked=True,
        account_turnover_monthly_minor=None if pd.isna(turnover) else int(turnover),
    )


def main() -> None:
    from risk.metrics import gini
    from risk.rules_v1 import assess_risk

    parser = argparse.ArgumentParser(description="Генератор синтетических компаний с меткой дефолта")
    parser.add_argument("--n", type=int, default=20_000)
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--out", type=Path, default=Path(__file__).parent / "generated" / "companies.csv")
    args = parser.parse_args()

    df = generate_companies(args.n, args.seed)
    args.out.parent.mkdir(parents=True, exist_ok=True)
    df.to_csv(args.out, index=False)

    scores = np.array([assess_risk(to_application(row)).score for _, row in df.iterrows()])
    print(f"Сохранено {len(df)} компаний в {args.out}")
    print(f"Доля дефолтов: {df['default_12m'].mean():.1%}")
    print(f"Gini rules-v1: {gini(df['default_12m'], -scores):.3f}")
    print(f"Gini истинной PD (потолок): {gini(df['default_12m'], df['true_pd']):.3f}")


if __name__ == "__main__":
    main()
