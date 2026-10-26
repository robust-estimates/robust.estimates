package usecase

import (
	"context"
	"errors"
	"testing"

	"robust/application-service/internal/domain"
)

type externalFake struct {
	registry domain.RegistryData
	report   domain.CreditReport
	fail     bool
}

func (f externalFake) Registry(context.Context, string) (domain.RegistryData, error) {
	if f.fail {
		return domain.RegistryData{}, errors.New("down")
	}
	return f.registry, nil
}

func (f externalFake) CreditReport(context.Context, string) (domain.CreditReport, error) {
	if f.fail {
		return domain.CreditReport{}, errors.New("down")
	}
	return f.report, nil
}

func (f externalFake) Account(context.Context, string) (domain.AccountData, error) {
	if f.fail {
		return domain.AccountData{}, errors.New("down")
	}
	return domain.AccountData{Found: true, TurnoverMonthlyMinor: 100}, nil
}

func clean() externalFake {
	history := &domain.CreditHistory{Found: true, ActiveLoans: 1}
	return externalFake{registry: domain.RegistryData{Status: "ACTIVE", OkvedSection: "G"}, report: domain.CreditReport{Company: history, Owner: history}}
}

func TestCollectPassesExternalDataToScoring(t *testing.T) {
	s, r := setup("SCORING")
	s.External = clean()
	features, checks := s.collect(context.Background(), r.apps["app"])
	if features.Registry == nil || features.CreditReport == nil || features.Account == nil {
		t.Fatal("external data was not passed to scoring")
	}
	if !checksPassed(checks) {
		t.Fatalf("clean company failed checks: %+v", checks)
	}
}

func TestNegativeFindingsFailChecks(t *testing.T) {
	s, r := setup("SCORING")
	fake := clean()
	fake.registry.TaxDebtMinor = 12000000
	fake.report.Owner = &domain.CreditHistory{Found: true, OverdueCount12m: 3, MaxOverdueDays24m: 75}
	s.External = fake
	_, checks := s.collect(context.Background(), r.apps["app"])
	if checks[0].Status != "FAILED" || checks[0].Summary != "Компания действует, ОКВЭД G; налоговая задолженность 120 000 ₽" {
		t.Fatalf("tax debt not reported: %+v", checks[0])
	}
	if checks[1].Status != "FAILED" {
		t.Fatalf("overdues not reported: %+v", checks[1])
	}
}

func TestUnavailableSourcesDoNotInventData(t *testing.T) {
	s, r := setup("SCORING")
	s.External = externalFake{fail: true}
	features, checks := s.collect(context.Background(), r.apps["app"])
	if features.Registry != nil || features.CreditReport != nil || features.Account != nil {
		t.Fatal("data invented for unavailable source")
	}
	for _, c := range checks {
		if c.Status != "UNAVAILABLE" {
			t.Fatalf("unavailable source reported as %s", c.Status)
		}
	}
}

func TestCreditReportNeedsConsent(t *testing.T) {
	s, r := setup("SCORING")
	s.External = clean()
	a := r.apps["app"]
	a.Consents.CreditHistoryCheck = false
	features, checks := s.collect(context.Background(), a)
	if features.CreditReport != nil || checks[1].Status != "UNAVAILABLE" {
		t.Fatal("credit history requested without consent")
	}
}

func TestLiquidatedCompanyIsRejected(t *testing.T) {
	_, r := setup("SCORING")
	a := r.apps["app"]
	features := domain.Features{Registry: &domain.RegistryData{Status: "LIQUIDATED"}}
	Reject(&a, domain.Scoring{Score: 800, RiskLevel: "LOW"}, stopReason(features))
	if a.Status != "REJECTED" || a.Decision.Outcome != "REJECTED" || len(a.Decision.Offers) != 0 {
		t.Fatalf("liquidated company not rejected: %s", a.Status)
	}
}

func TestFailedCheckPreventsAutomaticApproval(t *testing.T) {
	_, r := setup("SCORING")
	a := r.apps["app"]
	Finish(&a, domain.Scoring{Score: 800, RiskLevel: "LOW", Checks: []domain.Check{{Source: "FNS_MOCK", Status: "FAILED"}}})
	if a.Status != "MANUAL_CHECK" || a.Decision != nil {
		t.Fatalf("failed check approved automatically: %s", a.Status)
	}
}

func TestRubles(t *testing.T) {
	for minor, want := range map[int64]string{0: "0 ₽", 99900: "999 ₽", 12000000: "120 000 ₽", 123456789000: "1 234 567 890 ₽"} {
		if got := Rubles(minor); got != want {
			t.Fatalf("Rubles(%d) = %q, want %q", minor, got, want)
		}
	}
}
