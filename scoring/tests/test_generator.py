import numpy as np
import pandas as pd
import pytest

from data.generator import (MAX_AMOUNT_MINOR, MIN_AMOUNT_MINOR, TARGET_DEFAULT_RATE,
                            generate_companies, make_inn, to_application)
from risk.metrics import gini
from risk.rules_v1 import assess_risk


@pytest.fixture(scope="module")
def companies() -> pd.DataFrame:
    return generate_companies(5_000, seed=7)


def _inn_is_valid(inn: str) -> bool:
    def check(digits, weights):
        return str(sum(int(d) * w for d, w in zip(digits, weights)) % 11 % 10)

    if len(inn) == 10:
        return inn[9] == check(inn, [2, 4, 10, 3, 5, 9, 4, 6, 8])
    return (inn[10] == check(inn, [7, 2, 4, 10, 3, 5, 9, 4, 6, 8])
            and inn[11] == check(inn, [3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8]))


def test_same_seed_gives_same_data():
    pd.testing.assert_frame_equal(generate_companies(500, seed=1), generate_companies(500, seed=1))


def test_default_rate_matches_target(companies):
    assert companies["default_12m"].mean() == pytest.approx(TARGET_DEFAULT_RATE, abs=0.015)


def test_inns_are_unique_valid_and_synthetic(companies):
    inns = companies["inn"]
    assert inns.is_unique
    assert inns.str.startswith("00").all()
    assert (inns.str.len() == np.where(companies["legal_form"] == "OOO", 10, 12)).all()
    assert inns.map(_inn_is_valid).all()


def test_known_inn_checksum():
    # Контрольные цифры считаются по алгоритму ФНС
    assert _inn_is_valid(make_inn(123, "OOO"))
    assert _inn_is_valid(make_inn(123, "IP"))


def test_values_within_api_limits(companies):
    assert companies["requested_amount_minor"].between(MIN_AMOUNT_MINOR, MAX_AMOUNT_MINOR).all()
    assert companies["requested_term_months"].between(1, 60).all()
    assert (companies["business_age_months"] >= 0).all()
    assert (companies[["monthly_revenue_minor", "monthly_expenses_minor",
                       "monthly_debt_payments_minor", "tax_debt_minor"]] >= 0).all().all()


def test_ip_has_no_separate_company_history(companies):
    ip = companies[companies["legal_form"] == "IP"]
    assert ip["company_overdues_30_plus_12m"].isna().all()


def test_risk_factors_point_the_right_way(companies):
    has_overdue = companies["owner_overdues_30_plus_12m"].fillna(0) > 0
    rate = companies.groupby(has_overdue)["default_12m"].mean()
    assert rate[True] > rate[False]


def test_rules_v1_separates_defaults(companies):
    scores = np.array([assess_risk(to_application(row)).score for _, row in companies.iterrows()])
    rules_gini = gini(companies["default_12m"], -scores)
    ceiling = gini(companies["default_12m"], companies["true_pd"])
    assert 0.3 < rules_gini < ceiling
