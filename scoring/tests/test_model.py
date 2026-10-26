import unittest
from credit.v1 import credit_pb2 as pb
from app.model import RuleModel


def clean_history():
    return pb.CreditHistory(found=True, overdue_count_12m=0, max_overdue_days_24m=0, active_loans=1)


class ModelTests(unittest.TestCase):
    def test_stable_business(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=36, monthly_revenue_minor=300000000, monthly_expenses_minor=100000000))
        self.assertEqual(result.level, "LOW")
        self.assertGreaterEqual(result.score, 700)
        self.assertTrue(any("Возраст бизнеса" in factor for factor in result.factors))

    def test_unrequested_sources_are_not_penalized(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=36, monthly_revenue_minor=300000000, monthly_expenses_minor=100000000))
        self.assertFalse(any("кредитной истории" in factor for factor in result.factors))

    def test_credit_bureau_overdues_increase_risk(self):
        base = dict(business_age_months=36, monthly_revenue_minor=300000000, monthly_expenses_minor=100000000)
        clean = RuleModel().evaluate(pb.Features(**base, credit_report=pb.CreditReport(company=clean_history(), owner=clean_history())))
        bad_owner = pb.CreditHistory(found=True, overdue_count_12m=3, max_overdue_days_24m=75, active_loans=4)
        bad = RuleModel().evaluate(pb.Features(**base, credit_report=pb.CreditReport(company=clean_history(), owner=bad_owner)))
        self.assertLess(bad.score, clean.score)
        self.assertTrue(any("Просрочки 30+ дней у собственника за год: 3" in factor for factor in bad.factors))

    def test_no_history_found_is_reported(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=36, monthly_revenue_minor=300000000, monthly_expenses_minor=100000000,
                                                  registry=pb.RegistryData(status="ACTIVE", legal_form="IP"),
                                                  credit_report=pb.CreditReport(owner=pb.CreditHistory(found=False))))
        self.assertTrue(any("Нет кредитной истории собственника" in factor for factor in result.factors))
        self.assertFalse(any("компании" in factor for factor in result.factors))

    def test_registry_data_is_used(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=36, monthly_revenue_minor=300000000, monthly_expenses_minor=100000000,
                                                  registry=pb.RegistryData(status="ACTIVE", legal_form="OOO", okved_section="I", tax_debt_minor=12000000),
                                                  account=pb.AccountData(found=True, turnover_monthly_minor=250000000)))
        factors = " ".join(result.factors)
        self.assertIn("Налоговый долг 120 000 ₽", factors)
        self.assertIn("Отрасль: общепит и гостиницы", factors)
        self.assertIn("Есть обороты по счёту в банке", factors)

    def test_debt_load_increases_risk(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=6, monthly_revenue_minor=100000000, monthly_expenses_minor=90000000, monthly_debt_payments_minor=50000000))
        self.assertEqual(result.level, "HIGH")

    def test_zero_revenue(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=36))
        self.assertEqual(result.level, "HIGH")
        self.assertTrue(any("Нет выручки" in factor for factor in result.factors))

    def test_affordability_uses_requested_term_and_product_rate(self):
        product = pb.Product(id="unsecured", annual_rate_bps=1690, payment_type="ANNUITY")
        features = pb.Features(business_age_months=36, monthly_revenue_minor=30000000, monthly_expenses_minor=15000000, requested_amount_minor=500000000, term_months=12)
        result = RuleModel().evaluate(features, product)
        self.assertEqual(result.affordability.status, "CALCULATED")
        self.assertEqual(result.affordability.max_payment_minor, 12000000)
        self.assertLess(result.affordability.max_amount_minor, 500000000)

    def test_no_product_means_no_affordability(self):
        self.assertIsNone(RuleModel().evaluate(pb.Features(business_age_months=36, monthly_revenue_minor=1)).affordability)
