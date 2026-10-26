from dataclasses import dataclass
from decimal import ROUND_FLOOR, Decimal

MIN_DSCR = Decimal("1.20")
STRESS_FACTOR = Decimal("0.20")
MIN_STRESS_DSCR = Decimal("1.00")
AMOUNT_STEP_MINOR = 100_000

ESTIMATE_PAYMENT_TYPES = {"ANNUITY_ESTIMATE"}
SUPPORTED_PAYMENT_TYPES = {"ANNUITY"} | ESTIMATE_PAYMENT_TYPES


@dataclass(frozen=True)
class Affordability:
    status: str
    cfads_minor: int
    max_payment_minor: int
    max_amount_minor: int
    dscr: float
    reason: str


def _rubles(minor: int) -> str:
    return f"{minor // 100:,}".replace(",", " ") + " ₽"


def annuity_payment(amount_minor: int, term_months: int, annual_rate_bps: int) -> int:
    rate = Decimal(annual_rate_bps) / Decimal(120_000)
    if rate == 0:
        return int((Decimal(amount_minor) / term_months).to_integral_value())
    factor = (1 + rate) ** term_months
    return int((Decimal(amount_minor) * rate * factor / (factor - 1)).to_integral_value())


def max_amount(payment_minor: int, term_months: int, annual_rate_bps: int) -> int:
    rate = Decimal(annual_rate_bps) / Decimal(120_000)
    if rate == 0:
        amount = Decimal(payment_minor) * term_months
    else:
        factor = (1 + rate) ** term_months
        amount = Decimal(payment_minor) * (factor - 1) / (rate * factor)
    amount = int(amount.to_integral_value(rounding=ROUND_FLOOR)) // AMOUNT_STEP_MINOR * AMOUNT_STEP_MINOR
    while amount > 0 and annuity_payment(amount, term_months, annual_rate_bps) > payment_minor:
        amount -= AMOUNT_STEP_MINOR
    return max(amount, 0)


def _insufficient(reason: str) -> Affordability:
    return Affordability("INSUFFICIENT_DATA", 0, 0, 0, 0.0, reason)


def assess(revenue_minor: int, expenses_minor: int, debt_payments_minor: int, requested_amount_minor: int,
           term_months: int, annual_rate_bps: int, payment_type: str) -> Affordability:
    if revenue_minor <= 0:
        return _insufficient("Нет данных о выручке: платёжеспособность не рассчитана")
    if term_months < 1 or annual_rate_bps < 0:
        return _insufficient("Не заданы срок или ставка продукта")
    if payment_type not in SUPPORTED_PAYMENT_TYPES:
        return _insufficient("Для этого типа платежей нет модели расчёта")

    cfads = revenue_minor - expenses_minor
    base_capacity = int((Decimal(cfads) / MIN_DSCR).to_integral_value(rounding=ROUND_FLOOR))
    stress_capacity = int((Decimal(cfads) * (1 - STRESS_FACTOR) / MIN_STRESS_DSCR).to_integral_value(rounding=ROUND_FLOOR))
    max_payment = max(min(base_capacity, stress_capacity) - debt_payments_minor, 0)

    payment = annuity_payment(requested_amount_minor, term_months, annual_rate_bps) if requested_amount_minor > 0 else 0
    total_debt = debt_payments_minor + payment
    dscr = float(Decimal(cfads) / total_debt) if total_debt > 0 else 0.0

    if max_payment == 0:
        reason = (f"Свободного денежного потока {_rubles(max(cfads, 0))} в месяц не хватает на новый платёж "
                  f"с учётом стресс-сценария −20% и текущих платежей {_rubles(debt_payments_minor)}")
        return Affordability("CALCULATED", cfads, 0, 0, dscr, reason)

    reason = (f"Платёж по новому кредиту не должен превышать {_rubles(max_payment)} в месяц: "
              f"свободный денежный поток {_rubles(cfads)} с учётом стресс-сценария −20% "
              f"и текущих платежей {_rubles(debt_payments_minor)}")
    if payment_type in ESTIMATE_PAYMENT_TYPES:
        reason += "; для лимита — оценка при полном использовании"
    return Affordability("CALCULATED", cfads, max_payment, max_amount(max_payment, term_months, annual_rate_bps), dscr, reason)
