from decimal import Decimal

import pytest

from app.product_rules import (
    ApplicationData,
    AssessmentStatus,
    PaymentModel,
    ProductId,
    ProductTerms,
    annuity_payment,
    assess_all_products,
    assess_product,
)


def make_application(**overrides):
    data = {
        "business_age_months": 48,
        "monthly_revenue": Decimal("1000000"),
        "monthly_expenses": Decimal("600000"),
        "monthly_debt_payments": Decimal("50000"),
        "requested_amount": Decimal("1000000"),
        "requested_term_months": 36,
        "purpose": "WORKING_CAPITAL",
        "has_collateral": False,
    }
    data.update(overrides)
    return ApplicationData(**data)


def test_annuity_payment_zero_interest():
    assert annuity_payment(
        Decimal("120000"),
        Decimal("0"),
        12,
    ) == Decimal("10000")


def test_annuity_payment_positive_interest():
    payment = annuity_payment(
        Decimal("1000000"),
        Decimal("24"),
        12,
    )
    assert payment > Decimal("90000")
    assert payment < Decimal("100000")


def test_annuity_rejects_invalid_term():
    with pytest.raises(ValueError):
        annuity_payment(Decimal("1000"), Decimal("10"), 0)


def test_business_loan_can_be_assessed():
    app = make_application()
    terms = ProductTerms(
        product_id=ProductId.BUSINESS_LOAN_UNSECURED,
        min_amount=Decimal("100000"),
        max_amount=Decimal("5000000"),
        min_term_months=12,
        max_term_months=60,
        annual_rate_percent=Decimal("20"),
        payment_model=PaymentModel.ANNUITY,
    )

    result = assess_product(app, terms)

    assert result.status in {
        AssessmentStatus.AVAILABLE,
        AssessmentStatus.UNAVAILABLE,
    }
    assert result.estimated_monthly_payment is not None
    assert result.maximum_affordable_amount is not None


def test_secured_product_requires_collateral():
    app = make_application(has_collateral=False)
    terms = ProductTerms(
        product_id=ProductId.BUSINESS_LOAN_SECURED,
        requires_collateral=True,
        payment_model=PaymentModel.ANNUITY,
        annual_rate_percent=Decimal("20"),
        max_term_months=36,
    )

    result = assess_product(app, terms)

    assert result.status == AssessmentStatus.UNAVAILABLE


def test_missing_financial_data_is_not_treated_as_zero():
    app = make_application(monthly_revenue=None)
    terms = ProductTerms(
        product_id=ProductId.BUSINESS_LOAN_UNSECURED,
        annual_rate_percent=Decimal("20"),
        max_term_months=36,
        payment_model=PaymentModel.ANNUITY,
    )

    result = assess_product(app, terms)

    assert result.status == AssessmentStatus.INSUFFICIENT_DATA


def test_credit_line_does_not_use_annuity_automatically():
    app = make_application()
    terms = ProductTerms(
        product_id=ProductId.REVOLVING_LINE_UNSECURED,
        annual_rate_percent=Decimal("20"),
        max_term_months=36,
        payment_model=PaymentModel.CREDIT_LINE,
    )

    result = assess_product(app, terms)

    assert result.status == AssessmentStatus.MANUAL_REVIEW
    assert result.estimated_monthly_payment is None


def test_all_12_product_ids_are_supported():
    assert len(ProductId) == 12


def test_catalog_returns_one_result_per_product():
    app = make_application()
    catalog = [
        ProductTerms(
            product_id=product_id,
            payment_model=PaymentModel.MANUAL,
        )
        for product_id in ProductId
    ]

    results = assess_all_products(app, catalog)

    assert len(results) == 12
