from dataclasses import dataclass

from risk.rules_v1 import MODEL_VERSION, assess_risk
from schema import Application, CreditHistory


@dataclass(frozen=True)
class Risk:
    score: int
    level: str
    factors: tuple[str, ...]
    recommendation: str


class RuleModel:
    version = MODEL_VERSION

    def evaluate(self, features):
        overdue = bool(features.has_overdue_debt)
        application = Application(
            business_age_months=features.business_age_months,
            monthly_revenue_minor=features.monthly_revenue_minor,
            monthly_expenses_minor=features.monthly_expenses_minor,
            monthly_debt_payments_minor=features.monthly_debt_payments_minor,
            requested_amount_minor=features.requested_amount_minor,
            requested_term_months=features.term_months,
            company_credit_history=CreditHistory(
                overdues_30_plus_12m=1 if overdue else 0,
                max_overdue_days_24m=90 if overdue else 0,
                active_loans=1 if features.monthly_debt_payments_minor > 0 else 0,
            ),
        )
        result = assess_risk(application)
        level = result.risk_level.value
        recommendation = {"LOW": "APPROVE", "MEDIUM": "MANUAL_REVIEW", "HIGH": "REJECT"}[level]

        return Risk(result.score, level, tuple(factor.text for factor in result.factors), recommendation)
