"""Общий формат входных данных скоринга."""

from dataclasses import dataclass
from enum import StrEnum


class Purpose(StrEnum):
    WORKING_CAPITAL = "WORKING_CAPITAL"
    EQUIPMENT = "EQUIPMENT"
    BUSINESS_EXPANSION = "BUSINESS_EXPANSION"


class LegalForm(StrEnum):
    IP = "IP"
    OOO = "OOO"


class RiskLevel(StrEnum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"


@dataclass(frozen=True)
class CreditHistory:
    """Сводка из БКИ (мок) по компании или собственнику."""

    overdues_30_plus_12m: int  # число просрочек от 30 дней за последние 12 месяцев
    max_overdue_days_24m: int  # самая длинная просрочка за 24 месяца, дней
    active_loans: int


@dataclass(frozen=True)
class Application:
    # Анкета клиента (CompanyInput + Need)
    business_age_months: int
    monthly_revenue_minor: int
    monthly_expenses_minor: int
    monthly_debt_payments_minor: int
    requested_amount_minor: int | None = None
    requested_term_months: int | None = None
    purpose: Purpose | None = None
    has_collateral: bool | None = None

    # ЕГРЮЛ / ЕГРИП (мок)
    legal_form: LegalForm | None = None
    okved_section: str | None = None  # буква раздела ОКВЭД, например "G" — торговля
    region_code: str | None = None

    # ФНС (мок)
    tax_debt_minor: int | None = None

    # БКИ (мок)
    company_credit_history: CreditHistory | None = None
    owner_credit_history: CreditHistory | None = None

    # Внутренние данные банка: None, если клиент не обслуживается в банке
    account_turnover_monthly_minor: int | None = None
