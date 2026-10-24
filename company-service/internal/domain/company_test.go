package domain

import "testing"

func TestCalculationPreservesPrincipal(t *testing.T) {
	for _, test := range []struct {
		amount int64
		months int32
		rate   int32
	}{{100000000, 36, 1690}, {10000001, 12, 0}, {50000000000, 180, 1410}} {
		r := Calculate(test.amount, test.months, test.rate)
		var principal, total int64
		for _, p := range r.Schedule {
			if p.PrincipalMinor < 0 || p.BalanceMinor < 0 {
				t.Fatal("negative principal or balance")
			}

			principal += p.PrincipalMinor
			total += p.PaymentMinor
		}

		if principal != test.amount || total != r.TotalPaymentMinor || r.Schedule[len(r.Schedule)-1].BalanceMinor != 0 || r.OverpaymentMinor != total-test.amount {
			t.Fatalf("invalid schedule: %+v", r)
		}
	}
}

func TestCompanyValidationRejectsUnboundedAndNegativeValues(t *testing.T) {
	good := Input{INN: "7700000000", Name: "Синтетическая компания", BusinessAgeMonths: 36, MonthlyRevenueMinor: 100000000}
	if !ValidInput(good) {
		t.Fatal("valid company rejected")
	}

	bad := good
	bad.MonthlyExpensesMinor = -1
	if ValidInput(bad) {
		t.Fatal("negative expenses accepted")
	}

	bad = good
	bad.MonthlyRevenueMinor = 10000000000001
	if ValidInput(bad) {
		t.Fatal("unbounded revenue accepted")
	}
}
