import os
from dataclasses import dataclass

from affordability.dscr import Affordability, assess
from risk import model_v1, rules_v1
from schema import Application, CreditHistory, LegalForm


@dataclass(frozen=True)
class Risk:
    score: int
    level: str
    factors: tuple[str, ...]
    recommendation: str
    affordability: Affordability | None = None


def _present(message, field: str):
    return getattr(message, field) if message is not None and message.HasField(field) else None


def _history(report, field: str) -> CreditHistory | None:
    history = _present(report, field)
    if history is None or not history.found:
        return None
    return CreditHistory(
        overdues_30_plus_12m=history.overdue_count_12m,
        max_overdue_days_24m=history.max_overdue_days_24m,
        active_loans=history.active_loans,
    )


MODELS = {rules_v1.MODEL_VERSION: rules_v1.assess_risk, "model-v1": model_v1.assess_risk}


class RiskModel:
    def __init__(self, version: str | None = None):
        self.version = version or os.getenv("SCORING_MODEL", "model-v1")
        if self.version not in MODELS:
            raise ValueError(f"Unknown scoring model {self.version}")
        self.assess = MODELS[self.version]

    def evaluate(self, features, product=None):
        registry = _present(features, "registry")
        report = _present(features, "credit_report")
        account = _present(features, "account")

        application = Application(
            business_age_months=features.business_age_months,
            monthly_revenue_minor=features.monthly_revenue_minor,
            monthly_expenses_minor=features.monthly_expenses_minor,
            monthly_debt_payments_minor=features.monthly_debt_payments_minor,
            requested_amount_minor=features.requested_amount_minor,
            requested_term_months=features.term_months,
            legal_form=LegalForm(registry.legal_form) if registry and registry.legal_form in LegalForm else None,
            okved_section=registry.okved_section or None if registry else None,
            tax_debt_minor=registry.tax_debt_minor if registry else None,
            company_credit_history=_history(report, "company"),
            owner_credit_history=_history(report, "owner"),
            credit_bureau_checked=report is not None,
            account_turnover_monthly_minor=account.turnover_monthly_minor if account and account.found else None,
        )
        result = self.assess(application)
        level = result.risk_level.value
        recommendation = {"LOW": "APPROVE", "MEDIUM": "MANUAL_REVIEW", "HIGH": "REJECT"}[level]

        affordability = None
        if product is not None:
            affordability = assess(features.monthly_revenue_minor, features.monthly_expenses_minor,
                                   features.monthly_debt_payments_minor, features.requested_amount_minor,
                                   features.term_months, product.annual_rate_bps, product.payment_type)

        return Risk(result.score, level, tuple(factor.text for factor in result.factors), recommendation, affordability)
