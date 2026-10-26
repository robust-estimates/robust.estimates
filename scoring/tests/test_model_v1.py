import json

import pytest

from app.model import RiskModel
from risk import model_v1
from risk.train import MODEL_PATH, train
from schema import Application, CreditHistory, LegalForm, RiskLevel

CLEAN = CreditHistory(overdues_30_plus_12m=0, max_overdue_days_24m=0, active_loans=1)


@pytest.fixture(scope="module")
def model():
    return model_v1.load()


def make_app(**overrides) -> Application:
    fields = dict(business_age_months=40, monthly_revenue_minor=300_000_000, monthly_expenses_minor=150_000_000,
                  monthly_debt_payments_minor=0, legal_form=LegalForm.OOO, okved_section="G", tax_debt_minor=0,
                  company_credit_history=CLEAN, owner_credit_history=CLEAN, credit_bureau_checked=True,
                  account_turnover_monthly_minor=250_000_000)
    return Application(**(fields | overrides))


def test_committed_model_is_reproducible(model):
    assert json.loads(json.dumps(train())) == json.loads(MODEL_PATH.read_text(encoding="utf-8"))


def test_coefficients_have_expected_sign(model):
    assert all(value < 0 for value in model["metrics"]["coefficients"].values())


def test_beats_baseline_on_holdout(model):
    metrics = model["metrics"]
    assert metrics["gini_rules_v1_test"] < metrics["gini_test"] <= metrics["gini_ceiling_test"]


def test_levels_order_default_rates(model):
    levels = model["metrics"]["levels_test"]
    assert levels["LOW"]["default_rate"] < levels["MEDIUM"]["default_rate"] < levels["HIGH"]["default_rate"]


def test_points_follow_risk_direction(model):
    points = model["points"]
    assert points["owner_history"]["90+"] < points["owner_history"]["30+"] < points["owner_history"]["clean"]
    assert points["debt_load"]["50+"] < points["debt_load"]["20-35"] < points["debt_load"]["0"]
    assert points["margin"]["<0"] < points["margin"]["5-15"] < points["margin"]["25+"]
    assert points["tax_debt"]["yes"] < points["tax_debt"]["no"]
    assert points["business_age"]["<12"] < points["business_age"]["60+"]


def test_e2e_profiles(model):
    assert model_v1.assess_risk(make_app()).risk_level == RiskLevel.LOW
    overdue = CreditHistory(overdues_30_plus_12m=3, max_overdue_days_24m=75, active_loans=4)
    assert model_v1.assess_risk(make_app(business_age_months=6, owner_credit_history=overdue, okved_section="F")).risk_level != RiskLevel.LOW


def test_unrequested_sources_are_neutral():
    risk = model_v1.assess_risk(make_app(credit_bureau_checked=False, company_credit_history=None, owner_credit_history=None,
                                         tax_debt_minor=None, okved_section=None))
    assert not any(f.feature in ("company_history", "owner_history", "tax_debt", "industry") for f in risk.factors)


def test_zero_revenue_is_high_risk():
    risk = model_v1.assess_risk(make_app(monthly_revenue_minor=0, monthly_expenses_minor=0))
    assert risk.risk_level == RiskLevel.HIGH
    assert any(f.text.startswith("Нет выручки") for f in risk.factors)


def test_factor_texts_and_api_format():
    risk = model_v1.assess_risk(make_app(tax_debt_minor=12_000_000))
    assert any(f.text.startswith("Налоговый долг 120 000 ₽: −") for f in risk.factors)
    api = risk.to_api()
    assert api["modelVersion"] == "model-v1" and 0 <= api["score"] <= 1000


def test_model_selection(monkeypatch):
    assert RiskModel().version == "model-v1"
    monkeypatch.setenv("SCORING_MODEL", "rules-v1")
    assert RiskModel().version == "rules-v1"
    with pytest.raises(ValueError):
        RiskModel("unknown")
