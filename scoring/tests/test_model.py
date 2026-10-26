import unittest
from credit.v1 import credit_pb2 as pb
from app.model import RuleModel


class ModelTests(unittest.TestCase):
    def test_stable_business(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=36, monthly_revenue_minor=300000000, monthly_expenses_minor=100000000))
        self.assertEqual(result.level, "LOW")
        self.assertGreaterEqual(result.score, 700)
        self.assertTrue(any("Нет кредитной истории собственника" in factor for factor in result.factors))
        self.assertTrue(any("Возраст бизнеса" in factor for factor in result.factors))

    def test_debt_and_overdue_increase_risk(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=6, monthly_revenue_minor=100000000, monthly_expenses_minor=90000000, monthly_debt_payments_minor=50000000, has_overdue_debt=True))
        self.assertEqual(result.level, "HIGH")
        self.assertGreaterEqual(result.score, 0)

    def test_zero_revenue(self):
        result = RuleModel().evaluate(pb.Features(business_age_months=36))
        self.assertEqual(result.level, "HIGH")
        self.assertTrue(any("Нет выручки" in factor for factor in result.factors))
