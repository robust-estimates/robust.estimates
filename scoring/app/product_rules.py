from dataclasses import dataclass
from decimal import Decimal
from enum import StrEnum


class ProductId(StrEnum):
    BUSINESS_LOAN_UNSECURED = "business_loan_unsecured"
    BUSINESS_LOAN_SECURED = "business_loan_secured"
    REVOLVING_LINE_UNSECURED = "revolving_line_unsecured"
    REVOLVING_LINE_SECURED = "revolving_line_secured"
    NON_REVOLVING_LINE_SECURED = "non_revolving_line_secured"
    MONEY_FOR_BUSINESS = "money_for_business"
    BUSINESS_MORTGAGE = "business_mortgage"
    OVERDRAFT = "overdraft"
    SELLER_CREDIT = "seller_credit"
    REFINANCING = "refinancing"
    STATE_SUPPORTED = "state_supported"
    BUSINESS_ONE_CREDIT = "business_one_credit"


class PaymentModel(StrEnum):
    ANNUITY = "ANNUITY"
    CREDIT_LINE = "CREDIT_LINE"
    OVERDRAFT = "OVERDRAFT"
    SELLER = "SELLER"
    REFINANCING = "REFINANCING"
    STATE_PROGRAM = "STATE_PROGRAM"
    CARD = "CARD"
    MANUAL = "MANUAL"


class AssessmentStatus(StrEnum):
    AVAILABLE = "AVAILABLE"
    UNAVAILABLE = "UNAVAILABLE"
    INSUFFICIENT_DATA = "INSUFFICIENT_DATA"
    MANUAL_REVIEW = "MANUAL_REVIEW"


@dataclass(frozen=True)
class ProductTerms:


    product_id: ProductId
    min_amount: Decimal | None = None
    max_amount: Decimal | None = None
    min_term_months: int | None = None
    max_term_months: int | None = None
    annual_rate_percent: Decimal | None = None
    min_business_age_months: int | None = None
    requires_collateral: bool = False
    payment_model: PaymentModel = PaymentModel.MANUAL
    allowed_purposes: tuple[str, ...] = ()


@dataclass(frozen=True)
class ApplicationData:


    business_age_months: int | None
    monthly_revenue: Decimal | None
    monthly_expenses: Decimal | None
    monthly_debt_payments: Decimal | None
    requested_amount: Decimal | None
    requested_term_months: int | None

    purpose: str | None = None
    has_collateral: bool | None = None


    collateral_value: Decimal | None = None
    monthly_marketplace_payouts: Decimal | None = None
    monthly_marketplace_returns: Decimal | None = None
    refinancing_balance: Decimal | None = None
    monthly_account_inflows: Decimal | None = None
    state_program_id: str | None = None
    state_program_valid: bool | None = None
    property_verified: bool | None = None
    seller_data_verified: bool | None = None
    card_tariff_available: bool | None = None


@dataclass(frozen=True)
class ProductAssessment:
    product_id: str
    status: AssessmentStatus
    reason: str
    estimated_monthly_payment: Decimal | None = None
    maximum_affordable_amount: Decimal | None = None


ZERO = Decimal("0")
ONE = Decimal("1")
MIN_DSCR = Decimal("1.20")
STRESS_FACTOR = Decimal("0.20")
MIN_STRESS_DSCR = Decimal("1.00")
COLLATERAL_LTV = Decimal("0.60")


def annuity_payment(
    amount: Decimal,
    annual_rate_percent: Decimal,
    term_months: int,
) -> Decimal:


    if amount < ZERO:
        raise ValueError("Сумма не может быть отрицательной")
    if annual_rate_percent < ZERO:
        raise ValueError("Ставка не может быть отрицательной")
    if term_months <= 0:
        raise ValueError("Срок должен быть положительным")

    monthly_rate = annual_rate_percent / Decimal("1200")

    if monthly_rate == ZERO:
        return amount / Decimal(term_months)

    factor = (ONE + monthly_rate) ** term_months
    return amount * monthly_rate * factor / (factor - ONE)


def _missing_financial_data(app: ApplicationData) -> bool:
    return any(
        value is None
        for value in (
            app.monthly_revenue,
            app.monthly_expenses,
            app.monthly_debt_payments,
        )
    )


def calculate_capacity(
    app: ApplicationData,
    terms: ProductTerms,
) -> tuple[Decimal | None, str]:


    if _missing_financial_data(app):
        return None, "Недостаточно финансовых данных"

    revenue = app.monthly_revenue
    expenses = app.monthly_expenses
    debts = app.monthly_debt_payments

    if revenue is None or expenses is None or debts is None:
        return None, "Недостаточно финансовых данных"

    if revenue < ZERO or expenses < ZERO or debts < ZERO:
        return None, "Финансовые значения не могут быть отрицательными"

    if terms.annual_rate_percent is None:
        return None, "Не задана ставка продукта"

    if terms.max_term_months is None or terms.max_term_months <= 0:
        return None, "Не задан допустимый срок продукта"



    cfads = revenue - expenses
    stressed_cfads = cfads * (ONE - STRESS_FACTOR)

    available_payment = (
        cfads / MIN_DSCR - debts
    )
    stressed_available_payment = (
        stressed_cfads / MIN_STRESS_DSCR - debts
    )

    payment_capacity = min(
        available_payment,
        stressed_available_payment,
    )

    if payment_capacity <= ZERO:
        return ZERO, "Денежного потока недостаточно для нового платежа"


    rate = terms.annual_rate_percent / Decimal("1200")
    months = terms.max_term_months

    if rate == ZERO:
        amount = payment_capacity * months
    else:
        factor = (ONE + rate) ** months
        amount = payment_capacity * (factor - ONE) / (rate * factor)

    if terms.max_amount is not None:
        amount = min(amount, terms.max_amount)

    return max(ZERO, amount), "Сумма рассчитана по DSCR и стресс-сценарию"


def assess_product(
    app: ApplicationData,
    terms: ProductTerms,
) -> ProductAssessment:


    pid = terms.product_id

    def result(
        status: AssessmentStatus,
        reason: str,
        payment: Decimal | None = None,
        capacity: Decimal | None = None,
    ) -> ProductAssessment:
        return ProductAssessment(
            product_id=pid.value,
            status=status,
            reason=reason,
            estimated_monthly_payment=payment,
            maximum_affordable_amount=capacity,
        )


    if (
        terms.min_business_age_months is not None
        and app.business_age_months is not None
        and app.business_age_months < terms.min_business_age_months
    ):
        return result(
            AssessmentStatus.UNAVAILABLE,
            "Возраст бизнеса меньше минимального требования продукта",
        )

    if (
        terms.min_business_age_months is not None
        and app.business_age_months is None
    ):
        return result(
            AssessmentStatus.INSUFFICIENT_DATA,
            "Неизвестен возраст бизнеса",
        )

    if terms.allowed_purposes:
        if app.purpose is None:
            return result(
                AssessmentStatus.INSUFFICIENT_DATA,
                "Не указана цель финансирования",
            )
        if app.purpose not in terms.allowed_purposes:
            return result(
                AssessmentStatus.UNAVAILABLE,
                "Продукт не предназначен для указанной цели",
            )

    if terms.requires_collateral and app.has_collateral is not True:
        if app.has_collateral is None:
            return result(
                AssessmentStatus.INSUFFICIENT_DATA,
                "Нужно уточнить наличие залога",
            )
        return result(
            AssessmentStatus.UNAVAILABLE,
            "Для продукта требуется залог",
        )


    if pid == ProductId.MONEY_FOR_BUSINESS:
        if app.has_collateral is not True or app.collateral_value is None:
            return result(
                AssessmentStatus.INSUFFICIENT_DATA,
                "Нужны сведения об имуществе и его оценке",
            )

    elif pid == ProductId.BUSINESS_MORTGAGE:
        if app.property_verified is not True:
            return result(
                AssessmentStatus.MANUAL_REVIEW,
                "Нужно проверить объект недвижимости и условия сделки",
            )

    elif pid == ProductId.OVERDRAFT:
        if app.monthly_account_inflows is None:
            return result(
                AssessmentStatus.INSUFFICIENT_DATA,
                "Нужна история поступлений по счёту",
            )

    elif pid == ProductId.SELLER_CREDIT:
        if (
            app.monthly_marketplace_payouts is None
            or app.monthly_marketplace_returns is None
            or app.seller_data_verified is not True
        ):
            return result(
                AssessmentStatus.INSUFFICIENT_DATA,
                "Нужны подтверждённые выплаты и данные о возвратах",
            )

    elif pid == ProductId.REFINANCING:
        if app.refinancing_balance is None:
            return result(
                AssessmentStatus.INSUFFICIENT_DATA,
                "Нужны остатки рефинансируемых обязательств",
            )

    elif pid == ProductId.STATE_SUPPORTED:
        if not app.state_program_id or app.state_program_valid is None:
            return result(
                AssessmentStatus.INSUFFICIENT_DATA,
                "Нужно указать программу господдержки и проверить её условия",
            )
        if not app.state_program_valid:
            return result(
                AssessmentStatus.UNAVAILABLE,
                "Программа не действует или её условия не выполнены",
            )

    elif pid == ProductId.BUSINESS_ONE_CREDIT:
        if app.card_tariff_available is not True:
            return result(
                AssessmentStatus.MANUAL_REVIEW,
                "Нужно проверить тариф и правила расчёта обязательств по карте",
            )


    if terms.payment_model != PaymentModel.ANNUITY:
        return result(
            AssessmentStatus.MANUAL_REVIEW,
            "Нужна отдельная подтверждённая модель расчёта платежей",
        )

    if terms.annual_rate_percent is None:
        return result(
            AssessmentStatus.INSUFFICIENT_DATA,
            "Не задана ставка продукта",
        )

    if app.requested_amount is None or app.requested_term_months is None:
        return result(
            AssessmentStatus.INSUFFICIENT_DATA,
            "Не указаны желаемые сумма и срок",
        )

    if terms.min_amount is not None and app.requested_amount < terms.min_amount:
        return result(
            AssessmentStatus.UNAVAILABLE,
            "Сумма меньше минимальной для продукта",
        )

    if terms.max_amount is not None and app.requested_amount > terms.max_amount:
        return result(
            AssessmentStatus.UNAVAILABLE,
            "Сумма выше максимальной для продукта",
        )

    if (
        terms.min_term_months is not None
        and app.requested_term_months < terms.min_term_months
    ):
        return result(
            AssessmentStatus.UNAVAILABLE,
            "Срок меньше минимального для продукта",
        )

    if (
        terms.max_term_months is not None
        and app.requested_term_months > terms.max_term_months
    ):
        return result(
            AssessmentStatus.UNAVAILABLE,
            "Срок превышает максимальный для продукта",
        )

    capacity, capacity_reason = calculate_capacity(app, terms)

    if capacity is None:
        return result(AssessmentStatus.INSUFFICIENT_DATA, capacity_reason)

    payment = annuity_payment(
        app.requested_amount,
        terms.annual_rate_percent,
        app.requested_term_months,
    )

    if app.requested_amount > capacity:
        return result(
            AssessmentStatus.UNAVAILABLE,
            "Запрошенная сумма превышает расчётную платёжеспособность",
            payment=payment,
            capacity=capacity,
        )

    return result(
        AssessmentStatus.AVAILABLE,
        capacity_reason,
        payment=payment,
        capacity=capacity,
    )


def assess_all_products(
    app: ApplicationData,
    catalog: list[ProductTerms],
) -> list[ProductAssessment]:


    return [assess_product(app, terms) for terms in catalog]
