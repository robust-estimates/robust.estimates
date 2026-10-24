package usecase

import (
	"context"
	"sort"

	"robust/company-service/internal/domain"
	platform "robust/platform/core"
)

type Repository interface {
	Get(context.Context, string) (domain.Company, error)
	List(context.Context, string) ([]domain.Company, error)
	Change(context.Context, platform.Command, string, func(*domain.Company) error) (domain.Company, error)
	Products(context.Context) ([]domain.Product, error)
}

type Service struct{ Repo Repository }

func (s *Service) Get(ctx context.Context, id string) (domain.Company, error) {
	c, e := s.Repo.Get(ctx, id)
	if e != nil {
		return c, e
	}

	u := platform.Identity(ctx)
	if u.Role == "CLIENT" && c.OwnerID != u.Subject {
		return c, platform.NotFound()
	}

	return c, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Company, error) {
	return s.Repo.List(ctx, platform.Identity(ctx).Subject)
}

func (s *Service) Create(ctx context.Context, v domain.Input) (domain.Company, error) {
	if e := platform.Require(ctx, "CLIENT"); e != nil {
		return domain.Company{}, e
	}

	if !domain.ValidInput(v) {
		return domain.Company{}, platform.Invalid("Проверьте ИНН, название и финансовые показатели")
	}

	cmd, e := platform.GetCommand(ctx, "company.create", v, false)
	if e != nil {
		return domain.Company{}, e
	}

	return s.Repo.Change(ctx, cmd, "", func(c *domain.Company) error {
		*c = domain.Company{ID: platform.ID(), Data: v, OwnerID: cmd.User, VerificationStatus: "NOT_CHECKED", Version: 1}

		return nil
	})
}

func (s *Service) Update(ctx context.Context, id string, v domain.Input) (domain.Company, error) {
	if e := platform.Require(ctx, "CLIENT"); e != nil {
		return domain.Company{}, e
	}

	if !domain.ValidInput(v) {
		return domain.Company{}, platform.Invalid("Некорректные сведения компании")
	}

	cmd, e := platform.GetCommand(ctx, "company.update/"+id, v, false)
	if e != nil {
		return domain.Company{}, e
	}

	return s.Repo.Change(ctx, cmd, id, func(c *domain.Company) error {
		if c.OwnerID != cmd.User {
			return platform.NotFound()
		}

		if c.Data.INN != v.INN {
			return platform.Invalid("ИНН созданной компании неизменяем")
		}

		c.Data = v
		c.Version++
		c.VerificationStatus = "NOT_CHECKED"

		return nil
	})
}

func (s *Service) Verify(ctx context.Context, id string) (domain.Company, error) {
	if e := platform.Require(ctx, "CLIENT"); e != nil {
		return domain.Company{}, e
	}

	cmd, e := platform.GetCommand(ctx, "company.verify/"+id, id, false)
	if e != nil {
		return domain.Company{}, e
	}

	return s.Repo.Change(ctx, cmd, id, func(c *domain.Company) error {
		if c.OwnerID != cmd.User {
			return platform.NotFound()
		}

		c.Version++
		c.VerificationStatus = "VERIFIED"
		c.VerificationMessage = "Локальный мок ФНС: компания подтверждена; реальные реестры не запрашивались"

		return nil
	})
}

func (s *Service) Products(ctx context.Context) ([]domain.Product, error) {
	return s.Repo.Products(ctx)
}

func (s *Service) Product(ctx context.Context, id string) (domain.Product, error) {
	ps, e := s.Products(ctx)
	if e != nil {
		return domain.Product{}, e
	}

	for _, p := range ps {
		if p.ID == id {
			return p, nil
		}
	}

	return domain.Product{}, platform.NotFound()
}

func (s *Service) Calculate(ctx context.Context, in domain.CalculationInput) (domain.Calculation, error) {
	p, e := s.Product(ctx, in.ProductID)
	if e != nil {
		return domain.Calculation{}, e
	}

	if in.AmountMinor < p.MinAmountMinor || in.AmountMinor > p.MaxAmountMinor || in.TermMonths < p.MinTermMonths || in.TermMonths > p.MaxTermMonths {
		return domain.Calculation{}, platform.Invalid("Сумма или срок вне условий продукта")
	}

	v := domain.Calculate(in.AmountMinor, in.TermMonths, p.AnnualRateBps)
	if p.PaymentType != "ANNUITY" {
		v.Disclaimer = "Упрощённая оценка полного использования лимита; не фактический график кредитной линии/овердрафта"
	}

	return v, nil
}

func (s *Service) Recommend(ctx context.Context, in domain.RecommendationInput) ([]domain.Recommendation, error) {
	c, e := s.Get(ctx, in.CompanyID)
	if e != nil {
		return nil, e
	}

	if in.TermMonths < 1 || in.TermMonths > 180 {
		return nil, platform.Invalid("Срок от 1 до 180 месяцев")
	}

	if in.AmountMinor != nil && (*in.AmountMinor < 10000000 || *in.AmountMinor > 50000000000) {
		return nil, platform.Invalid("Некорректная сумма")
	}

	ps, e := s.Products(ctx)
	if e != nil {
		return nil, e
	}

	known := false
	for _, p := range ps {
		if domain.Supports(p, in.Purpose) {
			known = true
		}
	}

	if !known {
		return nil, platform.Invalid("Неизвестная цель финансирования")
	}

	result := []domain.Recommendation{}
	for _, p := range ps {
		rank := int32(1)
		availability := "BEST_MATCH"
		reason := "Соответствует цели, сроку и параметрам бизнеса"
		amount := c.Data.MonthlyRevenueMinor * 2
		if in.AmountMinor != nil {
			amount = *in.AmountMinor
		}

		if amount < p.MinAmountMinor {
			amount = p.MinAmountMinor
		}

		if amount > p.MaxAmountMinor {
			amount = p.MaxAmountMinor
		}

		if !domain.Supports(p, in.Purpose) {
			rank = 2
			availability = "ALTERNATIVE"
			reason = "Продукт доступен, но имеет другое целевое назначение"
		}

		if c.Data.BusinessAgeMonths < 12 || in.TermMonths < p.MinTermMonths || in.TermMonths > p.MaxTermMonths || (in.AmountMinor != nil && (*in.AmountMinor < p.MinAmountMinor || *in.AmountMinor > p.MaxAmountMinor)) || (p.Collateral == "REQUIRED" && in.Collateral == "NONE") || (p.ID == "sellers" && !c.Data.MarketplaceSales) {
			rank = 3
			availability = "UNAVAILABLE"
			reason = "Параметры компании, суммы, срока или обеспечения не соответствуют продукту"
		}

		months := in.TermMonths
		if months < p.MinTermMonths {
			months = p.MinTermMonths
		}

		if months > p.MaxTermMonths {
			months = p.MaxTermMonths
		}

		calc := domain.Calculate(amount, months, p.AnnualRateBps)
		result = append(result, domain.Recommendation{Product: p, Reason: reason, EstimatedMonthlyPaymentMinor: calc.MonthlyPaymentMinor, Availability: availability, Rank: rank})
	}

	sort.SliceStable(result, func(i, j int) bool { return result[i].Rank < result[j].Rank })
	return result, nil
}
