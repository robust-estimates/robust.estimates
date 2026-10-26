package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"robust/application-service/internal/domain"
	platform "robust/platform/core"
)

type ExternalSources interface {
	Registry(context.Context, string) (domain.RegistryData, error)
	CreditReport(context.Context, string) (domain.CreditReport, error)
	Account(context.Context, string) (domain.AccountData, error)
}

func (s *Service) collect(ctx context.Context, a domain.Aggregate) (domain.Features, []domain.Check) {
	f := a.CompanySnapshot.Data
	features := domain.Features{BusinessAgeMonths: f.BusinessAgeMonths, MonthlyRevenueMinor: f.MonthlyRevenueMinor, MonthlyExpensesMinor: f.MonthlyExpensesMinor, MonthlyDebtPaymentsMinor: f.MonthlyDebtPaymentsMinor, RequestedAmountMinor: a.Need.AmountMinor, TermMonths: a.Need.TermMonths}
	checks := []domain.Check{s.registryCheck(ctx, f.INN, &features), s.creditCheck(ctx, f.INN, a.Consents, &features)}

	if s.External != nil {
		if account, e := s.External.Account(ctx, f.INN); e == nil {
			features.Account = &account
		} else {
			slog.Warn("bank account data unavailable", "application_id", a.ID, "error", e)
		}
	}

	return features, checks
}

func (s *Service) registryCheck(ctx context.Context, inn string, features *domain.Features) domain.Check {
	check := domain.Check{Source: "FNS_MOCK", Status: "UNAVAILABLE", Summary: "ФНС не ответила, сведения о компании не получены", CheckedAt: platform.Now()}
	if s.External == nil {
		return check
	}

	r, e := s.External.Registry(ctx, inn)
	if e != nil {
		slog.Warn("registry request failed", "error", e)
		return check
	}

	features.Registry = &r
	switch {
	case r.Status != "ACTIVE":
		check.Status, check.Summary = "FAILED", "По данным ФНС компания не действует"
	case r.TaxDebtMinor > 0:
		check.Status, check.Summary = "FAILED", "Компания действует, ОКВЭД "+r.OkvedSection+"; налоговая задолженность "+Rubles(r.TaxDebtMinor)
	default:
		check.Status, check.Summary = "PASSED", "Компания действует, ОКВЭД "+r.OkvedSection+"; налоговой задолженности нет"
	}

	return check
}

func (s *Service) creditCheck(ctx context.Context, inn string, consent *domain.Consent, features *domain.Features) domain.Check {
	check := domain.Check{Source: "BKI_MOCK", Status: "UNAVAILABLE", Summary: "БКИ не ответило, кредитная история не получена", CheckedAt: platform.Now()}
	if consent == nil || !consent.CreditHistoryCheck {
		check.Summary = "Нет согласия на запрос кредитной истории"
		return check
	}

	if s.External == nil {
		return check
	}

	r, e := s.External.CreditReport(ctx, inn)
	if e != nil {
		slog.Warn("credit report request failed", "error", e)
		return check
	}

	features.CreditReport = &r
	parts := []string{}
	if r.Company != nil {
		parts = append(parts, "компания: "+describeHistory(r.Company))
	}
	if r.Owner != nil {
		parts = append(parts, "собственник: "+describeHistory(r.Owner))
	}

	check.Status, check.Summary = "PASSED", strings.Join(parts, "; ")
	if hasOverdues(r.Company) || hasOverdues(r.Owner) {
		check.Status = "FAILED"
	}

	return check
}

func hasOverdues(h *domain.CreditHistory) bool {
	return h != nil && h.Found && h.OverdueCount12m > 0
}

func describeHistory(h *domain.CreditHistory) string {
	switch {
	case !h.Found:
		return "кредитная история не найдена"
	case h.OverdueCount12m > 0:
		return fmt.Sprintf("просрочек от 30 дней за год: %d, самая длинная %d дн.", h.OverdueCount12m, h.MaxOverdueDays24m)
	default:
		return "просрочек нет"
	}
}

func Rubles(minor int64) string {
	digits := strconv.FormatInt(minor/100, 10)
	var b strings.Builder
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteRune(' ')
		}
		b.WriteRune(d)
	}

	return b.String() + " ₽"
}

func stopReason(f domain.Features) string {
	if f.Registry != nil && f.Registry.Status != "ACTIVE" {
		return "По данным ФНС компания не действует"
	}

	return ""
}

func checksPassed(checks []domain.Check) bool {
	for _, c := range checks {
		if c.Status != "PASSED" {
			return false
		}
	}

	return true
}
