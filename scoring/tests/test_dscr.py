import pytest

from affordability.dscr import AMOUNT_STEP_MINOR, annuity_payment, assess, max_amount


def test_short_term_cannot_exceed_payment_capacity():
    result = assess(100_000_000, 60_000_000, 5_000_000, 800_000_000, 12, 2000, "ANNUITY")
    assert result.max_payment_minor == 27_000_000
    assert result.max_amount_minor < 800_000_000
    assert annuity_payment(result.max_amount_minor, 12, 2000) <= result.max_payment_minor


def test_longer_term_gives_larger_amount():
    short = assess(100_000_000, 60_000_000, 0, 100_000_000, 12, 2000, "ANNUITY")
    long = assess(100_000_000, 60_000_000, 0, 100_000_000, 60, 2000, "ANNUITY")
    assert long.max_amount_minor > short.max_amount_minor


@pytest.mark.parametrize("payment,term,rate", [(27_000_000, 12, 2000), (1_234_567, 7, 1690), (500_000_00, 60, 0), (99_999, 3, 10000)])
def test_max_amount_is_largest_affordable_step(payment, term, rate):
    amount = max_amount(payment, term, rate)
    assert amount % AMOUNT_STEP_MINOR == 0
    assert annuity_payment(amount, term, rate) <= payment
    assert annuity_payment(amount + AMOUNT_STEP_MINOR, term, rate) > payment


def test_zero_rate():
    assert annuity_payment(120_000_00, 12, 0) == 10_000_00
    assert max_amount(10_000_00, 12, 0) == 120_000_00


def test_payment_matches_go_domain_payment():
    for args, go_value in {(100_000_000, 12, 1690): 9_115_730, (131_900_000, 12, 1690): 12_023_648,
                           (120_000, 12, 0): 10_000, (50_000_000, 60, 1390): 1_160_822}.items():
        assert annuity_payment(*args) == go_value


def test_stress_scenario_is_binding():
    result = assess(1_000_000_00, 600_000_00, 0, 100_000_00, 12, 1690, "ANNUITY")
    assert result.max_payment_minor == 320_000_00


def test_existing_debt_consumes_capacity():
    result = assess(1_000_000_00, 600_000_00, 400_000_00, 100_000_00, 12, 1690, "ANNUITY")
    assert result.status == "CALCULATED"
    assert result.max_amount_minor == 0
    assert "не хватает" in result.reason


def test_negative_cash_flow():
    result = assess(1_000_000_00, 1_200_000_00, 0, 100_000_00, 12, 1690, "ANNUITY")
    assert result.max_payment_minor == 0 and result.max_amount_minor == 0


def test_missing_revenue_is_not_calculated():
    assert assess(0, 0, 0, 100_000_00, 12, 1690, "ANNUITY").status == "INSUFFICIENT_DATA"


def test_unknown_payment_type_is_not_calculated():
    assert assess(1_000_000_00, 0, 0, 100_000_00, 12, 1690, "CARD").status == "INSUFFICIENT_DATA"


def test_credit_line_is_marked_as_estimate():
    assert "оценка при полном использовании" in assess(1_000_000_00, 0, 0, 100_000_00, 12, 1590, "ANNUITY_ESTIMATE").reason


def test_reason_explains_payment_limit():
    reason = assess(3_000_000_00, 1_500_000_00, 0, 1_000_000_00, 12, 1690, "ANNUITY").reason
    assert reason.startswith("Платёж по новому кредиту не должен превышать 1 200 000 ₽ в месяц")


def test_dscr_for_requested_amount():
    result = assess(1_000_000_00, 600_000_00, 0, 100_000_000, 12, 1690, "ANNUITY")
    assert result.dscr == pytest.approx(400_000_00 / 9_115_730)
