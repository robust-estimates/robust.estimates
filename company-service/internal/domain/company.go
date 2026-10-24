package domain

import (
	"math"
	"regexp"
	"strings"
)

type Input struct {
	INN                      string `json:"inn"`
	Name                     string `json:"name"`
	BusinessAgeMonths        int32  `json:"businessAgeMonths"`
	MonthlyRevenueMinor      int64  `json:"monthlyRevenueMinor"`
	MonthlyExpensesMinor     int64  `json:"monthlyExpensesMinor"`
	MonthlyDebtPaymentsMinor int64  `json:"monthlyDebtPaymentsMinor"`
	KPP                      string `json:"kpp"`
	OGRN                     string `json:"ogrn"`
	Industry                 string `json:"industry"`
	RegisteredAt             string `json:"registeredAt"`
	RepresentativeName       string `json:"representativeName"`
	RepresentativePosition   string `json:"representativePosition"`
	Phone                    string `json:"phone"`
	Email                    string `json:"email"`
	AnnualRevenueMinor       int64  `json:"annualRevenueMinor"`
	AnnualExpensesMinor      int64  `json:"annualExpensesMinor"`
	CashFlowMinor            int64  `json:"cashFlowMinor"`
	OutstandingDebtMinor     int64  `json:"outstandingDebtMinor"`
	HasOverdueDebt           bool   `json:"hasOverdueDebt"`
	MarketplaceSales         bool   `json:"marketplaceSales"`
}

type Company struct {
	ID                  string `json:"id"`
	Data                Input  `json:"data"`
	VerificationStatus  string `json:"verificationStatus"`
	VerificationMessage string `json:"verificationMessage"`
	OwnerID             string `json:"ownerId"`
	Version             int64  `json:"version"`
}

type Product struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Purposes              []string `json:"purposes"`
	MinAmountMinor        int64    `json:"minAmountMinor"`
	MaxAmountMinor        int64    `json:"maxAmountMinor"`
	MinTermMonths         int32    `json:"minTermMonths"`
	MaxTermMonths         int32    `json:"maxTermMonths"`
	AnnualRateBps         int32    `json:"annualRateBps"`
	RequiredDocumentTypes []string `json:"requiredDocumentTypes"`
	Description           string   `json:"description"`
	Collateral            string   `json:"collateral"`
	PaymentType           string   `json:"paymentType"`
	Tags                  []string `json:"tags"`
}

type RecommendationInput struct {
	CompanyID   string `json:"companyId"`
	AmountMinor *int64 `json:"amountMinor"`
	TermMonths  int32  `json:"termMonths"`
	Purpose     string `json:"purpose"`
	Collateral  string `json:"collateral"`
}

type Recommendation struct {
	Product                      Product `json:"product"`
	Reason                       string  `json:"reason"`
	EstimatedMonthlyPaymentMinor int64   `json:"estimatedMonthlyPaymentMinor"`
	Availability                 string  `json:"availability"`
	Rank                         int32   `json:"rank"`
}

type CalculationInput struct {
	ProductID   string `json:"productId"`
	AmountMinor int64  `json:"amountMinor"`
	TermMonths  int32  `json:"termMonths"`
}

type Payment struct {
	Month          int32 `json:"month"`
	PrincipalMinor int64 `json:"principalMinor"`
	InterestMinor  int64 `json:"interestMinor"`
	PaymentMinor   int64 `json:"paymentMinor"`
	BalanceMinor   int64 `json:"balanceMinor"`
}

type Calculation struct {
	MonthlyPaymentMinor int64     `json:"monthlyPaymentMinor"`
	TotalPaymentMinor   int64     `json:"totalPaymentMinor"`
	OverpaymentMinor    int64     `json:"overpaymentMinor"`
	Schedule            []Payment `json:"schedule"`
	Disclaimer          string    `json:"disclaimer"`
}

func ValidInput(v Input) bool {
	return regexp.MustCompile(`^(\d{10}|\d{12})$`).MatchString(v.INN) && strings.TrimSpace(v.Name) != "" && len(v.Name) <= 300 && v.BusinessAgeMonths >= 0 && v.BusinessAgeMonths <= 2400 && v.MonthlyRevenueMinor >= 0 && v.MonthlyRevenueMinor <= 10000000000000 && v.MonthlyExpensesMinor >= 0 && v.MonthlyExpensesMinor <= 10000000000000 && v.MonthlyDebtPaymentsMinor >= 0 && v.MonthlyDebtPaymentsMinor <= 10000000000000 && v.AnnualRevenueMinor >= 0 && v.AnnualRevenueMinor <= 100000000000000 && v.AnnualExpensesMinor >= 0 && v.AnnualExpensesMinor <= 100000000000000 && v.OutstandingDebtMinor >= 0 && v.OutstandingDebtMinor <= 10000000000000 && v.CashFlowMinor >= -100000000000000 && v.CashFlowMinor <= 100000000000000
}

func Supports(p Product, purpose string) bool {
	for _, v := range p.Purposes {
		if v == purpose {
			return true
		}
	}

	return false
}

func Calculate(amount int64, months int32, bps int32) Calculation {
	r := float64(bps) / 120000
	pay := float64(amount) / float64(months)
	if r > 0 {
		pay = float64(amount) * r / (1 - math.Pow(1+r, -float64(months)))
	}

	monthly := int64(math.Round(pay))
	balance := amount
	result := Calculation{MonthlyPaymentMinor: monthly, Schedule: []Payment{}, Disclaimer: "Оценочный аннуитетный расчёт без комиссий и страхования; условия определяются решением банка"}
	for i := int32(1); i <= months; i++ {
		interest := int64(math.Round(float64(balance) * r))
		principal := monthly - interest
		if i == months || principal > balance {
			principal = balance
		}

		payment := principal + interest
		balance -= principal
		result.TotalPaymentMinor += payment
		result.Schedule = append(result.Schedule, Payment{i, principal, interest, payment, balance})
	}

	result.OverpaymentMinor = result.TotalPaymentMinor - amount

	return result
}

func Products() []Product {
	return []Product{
		{"unsecured", "Бизнес-кредит без залога", []string{"WORKING_CAPITAL", "BUSINESS_EXPANSION", "OTHER"}, 10000000, 3000000000, 6, 60, 1690, []string{"FINANCIAL_REPORT", "BANK_STATEMENT"}, "На развитие и оборотные средства", "NONE", "ANNUITY", []string{"POPULAR", "UNSECURED"}},
		{"secured", "Бизнес-кредит с залогом", []string{"WORKING_CAPITAL", "BUSINESS_EXPANSION", "EQUIPMENT", "OTHER"}, 100000000, 50000000000, 12, 120, 1390, []string{"FINANCIAL_REPORT", "BANK_STATEMENT", "REGISTRATION"}, "Крупные суммы под обеспечение", "REQUIRED", "ANNUITY", []string{"SECURED"}},
		{"credit-line", "Кредитная линия", []string{"WORKING_CAPITAL", "BUSINESS_EXPANSION"}, 10000000, 10000000000, 6, 36, 1590, []string{"FINANCIAL_REPORT", "BANK_STATEMENT"}, "Возобновляемый лимит; калькулятор показывает оценку полного использования", "NONE", "ANNUITY_ESTIMATE", []string{"REVOLVING"}},
		{"overdraft", "Овердрафт", []string{"WORKING_CAPITAL"}, 10000000, 5000000000, 1, 12, 1750, []string{"BANK_STATEMENT"}, "Для краткосрочных кассовых разрывов", "NONE", "ANNUITY_ESTIMATE", []string{"SHORT_TERM"}},
		{"mortgage", "Бизнес-ипотека", []string{"REAL_ESTATE"}, 100000000, 50000000000, 12, 180, 1410, []string{"FINANCIAL_REPORT", "BANK_STATEMENT", "REGISTRATION"}, "Коммерческая недвижимость", "REQUIRED", "ANNUITY", []string{"REAL_ESTATE"}},
		{"refinance", "Рефинансирование", []string{"REFINANCING"}, 100000000, 30000000000, 12, 120, 1350, []string{"FINANCIAL_REPORT", "BANK_STATEMENT"}, "Объединение существующих обязательств", "NONE", "ANNUITY", []string{"REFINANCING"}},
		{"sellers", "Кредит для селлеров", []string{"SELLERS", "WORKING_CAPITAL"}, 10000000, 5000000000, 6, 36, 1650, []string{"BANK_STATEMENT"}, "Для продаж на маркетплейсах", "NONE", "ANNUITY", []string{"SELLERS"}},
		{"business-one", "Business One Credit", []string{"BUSINESS_EXPANSION", "EQUIPMENT"}, 100000000, 20000000000, 12, 60, 1490, []string{"FINANCIAL_REPORT", "BANK_STATEMENT"}, "Комплексное финансирование бизнеса", "NONE", "ANNUITY", []string{"POPULAR"}},
	}
}
