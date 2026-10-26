from risk.rules_v1 import MODEL_VERSION, assess_risk
from schema import Application, CreditHistory, LegalForm, RiskLevel

CLEAN_HISTORY = CreditHistory(overdues_30_plus_12m=0, max_overdue_days_24m=0, active_loans=1)


def make_app(**overrides) -> Application:
    fields = dict(
        business_age_months=60,
        monthly_revenue_minor=300_000_000,
        monthly_expenses_minor=210_000_000,
        monthly_debt_payments_minor=30_000_000,
        tax_debt_minor=0,
        company_credit_history=CLEAN_HISTORY,
        owner_credit_history=CLEAN_HISTORY,
        account_turnover_monthly_minor=250_000_000,
        legal_form=LegalForm.OOO,
        credit_bureau_checked=True,
    )
    return Application(**(fields | overrides))


def test_reliable_client_is_low_risk():
    assert assess_risk(make_app()).risk_level == RiskLevel.LOW


def test_overdues_and_high_debt_load_is_high_risk():
    bad_history = CreditHistory(overdues_30_plus_12m=3, max_overdue_days_24m=95, active_loans=4)
    app = make_app(monthly_debt_payments_minor=180_000_000, owner_credit_history=bad_history,
                   tax_debt_minor=5_000_000)
    assert assess_risk(app).risk_level == RiskLevel.HIGH


def test_missing_data_is_not_treated_as_clean_history():
    clean = assess_risk(make_app())
    unknown = assess_risk(make_app(company_credit_history=None, owner_credit_history=None,
                                   tax_debt_minor=None, account_turnover_monthly_minor=None))
    assert unknown.score < clean.score
    assert any("Нет кредитной истории" in f.text for f in unknown.factors)


def test_zero_revenue_does_not_crash():
    risk = assess_risk(make_app(monthly_revenue_minor=0, monthly_expenses_minor=0))
    assert any(f.text.startswith("Нет выручки") for f in risk.factors)


def test_score_is_clipped_to_api_range():
    terrible = CreditHistory(overdues_30_plus_12m=10, max_overdue_days_24m=365, active_loans=9)
    app = make_app(business_age_months=1, monthly_revenue_minor=0, monthly_expenses_minor=0,
                   company_credit_history=terrible, owner_credit_history=terrible, tax_debt_minor=1)
    assert assess_risk(app).score == 0


def test_factors_sorted_by_strength():
    points = [abs(f.points) for f in assess_risk(make_app()).factors]
    assert points == sorted(points, reverse=True)


def test_to_api_matches_scoring_schema():
    api = assess_risk(make_app()).to_api()
    assert set(api) == {"score", "riskLevel", "factors", "modelVersion"}
    assert isinstance(api["score"], int) and 0 <= api["score"] <= 1000
    assert api["riskLevel"] in {"LOW", "MEDIUM", "HIGH"}
    assert all(isinstance(f, str) for f in api["factors"])
    assert api["modelVersion"] == MODEL_VERSION


def test_unrequested_credit_history_is_not_penalized():
    risk = assess_risk(make_app(company_credit_history=None, owner_credit_history=None, credit_bureau_checked=False))
    assert not any("кредитной истории" in f.text for f in risk.factors)


def test_ip_has_no_company_history_factor():
    risk = assess_risk(make_app(legal_form=LegalForm.IP, company_credit_history=None))
    assert not any("компании" in f.text for f in risk.factors)


def test_industry_factor():
    assert any(f.text == "Отрасль: IT и связь: +30" for f in assess_risk(make_app(okved_section="J")).factors)
    assert not any(f.feature == "okved_section" for f in assess_risk(make_app(okved_section="G")).factors)
