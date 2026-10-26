"""rules-v1: балльная карта.

Бейзлайн и запасной вариант для модели риска. Баллы подобраны вручную по
направлениям факторов из ТЗ: долговая нагрузка, просрочки и нестабильные финансы
повышают риск, стабильный бизнес и чистая кредитная история понижают.
Чем больше баллов, тем ниже риск.
"""

from dataclasses import dataclass

from schema import Application, CreditHistory, RiskLevel

MODEL_VERSION = "rules-v1"

BASE_POINTS = 600
MIN_SCORE, MAX_SCORE = 0, 1000

LOW_FROM = 700
MEDIUM_FROM = 500


@dataclass(frozen=True)
class Factor:
    feature: str
    value: str
    points: int
    text: str  # для андеррайтера: «Долговая нагрузка 42%: −60»


@dataclass(frozen=True)
class Risk:
    score: int
    risk_level: RiskLevel
    factors: list[Factor]
    model_version: str = MODEL_VERSION

    def to_api(self) -> dict:
        """Поля Scoring из openapi.yaml, кроме checks: их собирает бэкенд по ответам моков."""
        return {
            "score": self.score,
            "riskLevel": self.risk_level.value,
            "factors": [f.text for f in self.factors],
            "modelVersion": self.model_version,
        }


def _points_text(points: int) -> str:
    return f"+{points}" if points > 0 else f"−{-points}" if points < 0 else "0"


def _factor(feature: str, value: str, points: int, label: str) -> Factor:
    return Factor(feature, value, points, f"{label}: {_points_text(points)}")


def _business_age(months: int) -> Factor:
    if months < 12:
        points = -120
    elif months < 36:
        points = -40
    else:
        points = 50
    return _factor("business_age_months", str(months), points, f"Возраст бизнеса {months} мес.")


def _revenue_based(app: Application) -> list[Factor]:
    revenue = app.monthly_revenue_minor
    if revenue <= 0:
        return [_factor("monthly_revenue_minor", "0", -200, "Нет выручки")]

    debt_load = app.monthly_debt_payments_minor / revenue
    if debt_load > 0.5:
        debt_points = -150
    elif debt_load > 0.3:
        debt_points = -60
    else:
        debt_points = 40

    margin = (revenue - app.monthly_expenses_minor) / revenue
    if margin < 0:
        margin_points = -150
    elif margin < 0.05:
        margin_points = -60
    elif margin > 0.2:
        margin_points = 40
    else:
        margin_points = 0

    return [
        _factor("debt_load", f"{debt_load:.2f}", debt_points, f"Долговая нагрузка {debt_load:.0%}"),
        _factor("margin", f"{margin:.2f}", margin_points, f"Маржа {margin:.0%}"),
    ]


def _credit_history(history: CreditHistory | None, feature: str, whose: str) -> Factor:
    if history is None:
        return _factor(feature, "нет данных", -20, f"Нет кредитной истории {whose}")
    if history.max_overdue_days_24m >= 90:
        return _factor(feature, f"{history.max_overdue_days_24m} дн.", -250,
                       f"Просрочка 90+ дней у {whose} за 2 года")
    if history.overdues_30_plus_12m > 0:
        return _factor(feature, str(history.overdues_30_plus_12m), -200,
                       f"Просрочки 30+ дней у {whose} за год: {history.overdues_30_plus_12m}")
    return _factor(feature, "0", 30, f"Нет просрочек у {whose}")


def _tax_debt(tax_debt_minor: int | None) -> Factor | None:
    if tax_debt_minor is None:
        return None
    if tax_debt_minor > 0:
        return _factor("tax_debt_minor", str(tax_debt_minor), -100,
                       f"Налоговый долг {tax_debt_minor / 100:,.0f} ₽".replace(",", " "))
    return None


def _bank_client(turnover_minor: int | None) -> Factor | None:
    if turnover_minor is not None and turnover_minor > 0:
        return _factor("account_turnover_monthly_minor", str(turnover_minor), 30,
                       "Есть обороты по счёту в банке")
    return None


def risk_level(score: int) -> RiskLevel:
    if score >= LOW_FROM:
        return RiskLevel.LOW
    if score >= MEDIUM_FROM:
        return RiskLevel.MEDIUM
    return RiskLevel.HIGH


def assess_risk(app: Application) -> Risk:
    factors = [
        _business_age(app.business_age_months),
        *_revenue_based(app),
        _credit_history(app.company_credit_history, "company_credit_history", "компании"),
        _credit_history(app.owner_credit_history, "owner_credit_history", "собственника"),
    ]
    factors += [f for f in (_tax_debt(app.tax_debt_minor), _bank_client(app.account_turnover_monthly_minor)) if f]

    score = BASE_POINTS + sum(f.points for f in factors)
    score = max(MIN_SCORE, min(MAX_SCORE, score))
    # Сначала самые сильные факторы: андеррайтеру важнее то, что сильнее всего сдвинуло оценку
    factors.sort(key=lambda f: abs(f.points), reverse=True)
    return Risk(score=score, risk_level=risk_level(score), factors=factors)
