from dataclasses import dataclass


@dataclass(frozen=True)
class Risk:
    score: int
    level: str
    factors: tuple[str, ...]
    recommendation: str


class RuleModel:
    version = "rules-v1"

    def evaluate(self, features):
        score = 800
        factors = []

        if features.business_age_months < 12:
            score -= 150
            factors.append("Короткий срок деятельности")

        if features.has_overdue_debt:
            score -= 250
            factors.append("Просрочки по синтетическим данным БКИ")

        if features.monthly_revenue_minor <= 0:
            score -= 300
            factors.append("Недостаточно данных об обороте")
        else:
            debt_ratio = features.monthly_debt_payments_minor / features.monthly_revenue_minor

            if debt_ratio > 0.4:
                score -= 200
                factors.append("Высокая текущая долговая нагрузка")

            cash = features.monthly_revenue_minor - features.monthly_expenses_minor - features.monthly_debt_payments_minor

            if cash <= 0:
                score -= 150
                factors.append("Отрицательный операционный остаток")

        score = max(0, min(1000, score))
        level = "LOW" if score >= 700 else "MEDIUM" if score >= 450 else "HIGH"
        recommendation = "APPROVE" if level == "LOW" else "MANUAL_REVIEW" if level == "MEDIUM" else "REJECT"

        return Risk(score, level, tuple(factors or ["Стабильный оборот и приемлемая долговая нагрузка"]), recommendation)
