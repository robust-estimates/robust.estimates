package domain

import (
	"math"
	"time"
)

type CompanyData struct {
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
	ID                  string      `json:"id"`
	Data                CompanyData `json:"data"`
	VerificationStatus  string      `json:"verificationStatus"`
	VerificationMessage string      `json:"verificationMessage"`
	OwnerID             string      `json:"ownerId"`
	Version             int64       `json:"version"`
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

type Need struct {
	CompanyID   string `json:"companyId"`
	AmountMinor int64  `json:"amountMinor"`
	TermMonths  int32  `json:"termMonths"`
	Purpose     string `json:"purpose"`
	Description string `json:"description"`
}

type Consent struct {
	PersonalDataProcessing bool `json:"personalDataProcessing"`
	CreditHistoryCheck     bool `json:"creditHistoryCheck"`
}

type DraftInput struct {
	Need      Need     `json:"need"`
	ProductID string   `json:"productId"`
	Consents  *Consent `json:"consents"`
}

type Patch struct {
	Need      *Need    `json:"need"`
	ProductID *string  `json:"productId"`
	Consents  *Consent `json:"consents"`
}

type Application struct {
	ID             string   `json:"id"`
	CompanyID      string   `json:"companyId"`
	ProductID      string   `json:"productId"`
	Need           Need     `json:"need"`
	Consents       *Consent `json:"consents"`
	Status         string   `json:"status"`
	Version        int64    `json:"version"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	NextAction     string   `json:"nextAction"`
	MissingFields  []string `json:"missingFields"`
	Number         string   `json:"number"`
	CompanyName    string   `json:"companyName"`
	INN            string   `json:"inn"`
	ProductName    string   `json:"productName"`
	Recommendation string   `json:"recommendation"`
	OwnerID        string   `json:"ownerId"`
}

type Document struct {
	ID            string `json:"id"`
	ApplicationID string `json:"applicationId"`
	Type          string `json:"type"`
	FileName      string `json:"fileName"`
	ContentType   string `json:"contentType"`
	SizeBytes     int64  `json:"sizeBytes"`
	Status        string `json:"status"`
	Checksum      string `json:"checksum"`
	RequestItemID string `json:"requestItemId"`
	ObjectKey     string `json:"objectKey"`
}

type Validation struct {
	Valid         bool         `json:"valid"`
	MissingFields []string     `json:"missingFields"`
	Errors        []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Offer struct {
	ID                  string `json:"id"`
	ProductID           string `json:"productId"`
	AmountMinor         int64  `json:"amountMinor"`
	TermMonths          int32  `json:"termMonths"`
	AnnualRateBps       int32  `json:"annualRateBps"`
	MonthlyPaymentMinor int64  `json:"monthlyPaymentMinor"`
	ExpiresAt           string `json:"expiresAt"`
	Collateral          string `json:"collateral"`
}

type Decision struct {
	ID        string  `json:"id"`
	Outcome   string  `json:"outcome"`
	Reason    string  `json:"reason"`
	Offers    []Offer `json:"offers"`
	DecidedAt string  `json:"decidedAt"`
}

type InformationItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Status   string `json:"status"`
}

type InformationRequest struct {
	ID        string            `json:"id"`
	Message   string            `json:"message"`
	Items     []InformationItem `json:"items"`
	CreatedAt string            `json:"createdAt"`
}

type ManualDecision struct {
	Outcome         string            `json:"outcome"`
	Reason          string            `json:"reason"`
	Offers          []Offer           `json:"offers"`
	RequestedFields []string          `json:"requestedFields"`
	Items           []InformationItem `json:"items"`
}

type Loan struct {
	ID            string `json:"id"`
	ApplicationID string `json:"applicationId"`
	OfferID       string `json:"offerId"`
	AmountMinor   int64  `json:"amountMinor"`
	IssuedAt      string `json:"issuedAt"`
	Simulated     bool   `json:"simulated"`
}

type Event struct {
	ID         string `json:"id"`
	At         string `json:"at"`
	FromStatus string `json:"fromStatus"`
	ToStatus   string `json:"toStatus"`
	Message    string `json:"message"`
	ActorRole  string `json:"actorRole"`
}

type Check struct {
	Source    string `json:"source"`
	Status    string `json:"status"`
	Summary   string `json:"summary"`
	CheckedAt string `json:"checkedAt"`
}

type Scoring struct {
	Score          int32    `json:"score"`
	RiskLevel      string   `json:"riskLevel"`
	Factors        []string `json:"factors"`
	ModelVersion   string   `json:"modelVersion"`
	Checks         []Check  `json:"checks"`
	Recommendation string   `json:"recommendation"`
}

type Features struct {
	BusinessAgeMonths        int32 `json:"businessAgeMonths"`
	MonthlyRevenueMinor      int64 `json:"monthlyRevenueMinor"`
	MonthlyExpensesMinor     int64 `json:"monthlyExpensesMinor"`
	MonthlyDebtPaymentsMinor int64 `json:"monthlyDebtPaymentsMinor"`
	RequestedAmountMinor     int64 `json:"requestedAmountMinor"`
	TermMonths               int32 `json:"termMonths"`
	HasOverdueDebt           bool  `json:"hasOverdueDebt"`
}

type ScoringInput struct {
	ScoringRequestID   string   `json:"scoringRequestId"`
	ApplicationID      string   `json:"applicationId"`
	ApplicationVersion int64    `json:"applicationVersion"`
	Features           Features `json:"features"`
}

type Aggregate struct {
	Application
	CompanySnapshot     Company              `json:"companySnapshot"`
	ProductSnapshot     Product              `json:"productSnapshot"`
	Documents           []Document           `json:"documents"`
	Decision            *Decision            `json:"decision"`
	Scoring             *Scoring             `json:"scoring"`
	Events              []Event              `json:"events"`
	InformationRequests []InformationRequest `json:"informationRequests"`
	Loan                *Loan                `json:"loan"`
	AcceptedOfferID     string               `json:"acceptedOfferId"`
	Submitted           bool                 `json:"submitted"`
}

type Result struct {
	Application *Application `json:"application,omitempty"`
	Document    *Document    `json:"document,omitempty"`
	Loan        *Loan        `json:"loan,omitempty"`
}

type PageInput struct {
	Offset         int32  `json:"offset"`
	Limit          int32  `json:"limit"`
	Status         string `json:"status"`
	Search         string `json:"search"`
	CompanyID      string `json:"companyId"`
	Recommendation string `json:"recommendation"`
	FromDate       string `json:"fromDate"`
	ToDate         string `json:"toDate"`
}

type Page struct {
	Items  []Application    `json:"items"`
	Offset int32            `json:"offset"`
	Limit  int32            `json:"limit"`
	Total  int64            `json:"total"`
	Counts map[string]int64 `json:"counts"`
}

type FeedItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	Reason        string `json:"reason"`
	CreatedAt     string `json:"createdAt"`
	ApplicationID string `json:"applicationId"`
}

type Job struct {
	ID            string
	ApplicationID string
	Revision      int64
	Attempts      int32
}

func CanEdit(s string) bool { return s == "IN_PROGRESS" || s == "NEEDS_INFO" }

func Next(s string) string {
	switch s {
	case "IN_PROGRESS":
		return "Заполните заявку и прикрепите документы"
	case "SCORING":
		return "Проводим проверки и оценку риска"
	case "MANUAL_CHECK":
		return "Заявка рассматривается сотрудником"
	case "NEEDS_INFO":
		return "Предоставьте запрошенные сведения и отправьте заявку повторно"
	case "APPROVED", "ALTERNATIVE":
		return "Ознакомьтесь с предложением и примите решение"
	case "ACCEPTED":
		return "Подтвердите моделирование выдачи"
	case "ISSUED":
		return "Кредит выдан в демонстрационном режиме"
	case "REJECTED":
		return "Получен отказ банка"
	case "DECLINED":
		return "Вы отказались от предложения"
	case "CANCELLED":
		return "Заявка отозвана"
	}

	return "Обработка"
}

func Supports(p Product, purpose string) bool {
	for _, x := range p.Purposes {
		if x == purpose {
			return true
		}
	}

	return false
}

func Payment(amount int64, months int32, bps int32) int64 {
	r := float64(bps) / 120000
	if r == 0 {
		return int64(math.Round(float64(amount) / float64(months)))
	}

	return int64(math.Round(float64(amount) * r / (1 - math.Pow(1+r, -float64(months)))))
}

func OfferValid(o Offer, p Product) bool {
	t, e := time.Parse(time.RFC3339Nano, o.ExpiresAt)

	return e == nil && t.After(time.Now()) && o.ProductID == p.ID && o.AmountMinor >= p.MinAmountMinor && o.AmountMinor <= p.MaxAmountMinor && o.TermMonths >= p.MinTermMonths && o.TermMonths <= p.MaxTermMonths && o.AnnualRateBps >= 0 && o.AnnualRateBps <= 10000
}
