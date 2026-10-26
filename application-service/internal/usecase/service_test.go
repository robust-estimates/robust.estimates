package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"robust/application-service/internal/domain"
	platform "robust/platform/core"
)

type receipt struct {
	Hash   string
	Result domain.Result
}

type memoryRepo struct {
	mu       sync.Mutex
	apps     map[string]domain.Aggregate
	receipts map[string]receipt
}

func clone[T any](v T) T { var out T; b, _ := json.Marshal(v); _ = json.Unmarshal(b, &out); return out }

func (r *memoryRepo) Get(ctx context.Context, id string) (domain.Aggregate, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.apps[id]
	if !ok {
		return a, platform.NotFound()
	}

	return clone(a), nil
}

func (r *memoryRepo) List(context.Context, string, domain.PageInput) (domain.Page, error) {
	return domain.Page{}, nil
}

func (r *memoryRepo) Apply(ctx context.Context, c platform.Command, id string, fn func(*domain.Aggregate) (domain.Result, error)) (domain.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := c.User + c.Operation + c.Key
	if old, ok := r.receipts[key]; ok {
		if old.Hash != c.Hash {
			return domain.Result{}, platform.Conflict("key conflict")
		}

		return clone(old.Result), nil
	}

	a := clone(r.apps[id])
	if id != "" && a.ID == "" {
		return domain.Result{}, platform.NotFound()
	}

	out, e := fn(&a)
	if e != nil {
		return domain.Result{}, e
	}

	a.Version++
	out.Application = &a.Application
	r.apps[a.ID] = clone(a)
	r.receipts[key] = receipt{c.Hash, clone(out)}

	return out, nil
}

func (r *memoryRepo) Claim(context.Context) (*domain.Job, error) { return nil, nil }

func (r *memoryRepo) Complete(context.Context, domain.Job, func(*domain.Aggregate) error, bool) error {
	return nil
}

func (r *memoryRepo) StageObject(context.Context, string) error { return nil }

func (r *memoryRepo) Garbage(context.Context) ([]string, error) { return nil, nil }

func (r *memoryRepo) RemoveObject(context.Context, string) error { return nil }

type companyFake struct{}

func (companyFake) Company(context.Context, string) (domain.Company, error) {
	return domain.Company{ID: "company", OwnerID: "client", VerificationStatus: "VERIFIED", Data: domain.CompanyData{MonthlyRevenueMinor: 300000000}}, nil
}

func (companyFake) Product(context.Context, string) (domain.Product, error) { return product(), nil }

func product() domain.Product {
	return domain.Product{ID: "unsecured", MinAmountMinor: 10000000, MaxAmountMinor: 3000000000, MinTermMonths: 6, MaxTermMonths: 60, AnnualRateBps: 1690, Purposes: []string{"WORKING_CAPITAL"}, RequiredDocumentTypes: []string{"BANK_STATEMENT"}}
}

func setup(status string) (*Service, *memoryRepo) {
	a := domain.Aggregate{Application: domain.Application{ID: "app", OwnerID: "client", CompanyID: "company", ProductID: "unsecured", Status: status, Version: 1, Need: domain.Need{CompanyID: "company", AmountMinor: 100000000, TermMonths: 12, Purpose: "WORKING_CAPITAL"}, Consents: &domain.Consent{PersonalDataProcessing: true, CreditHistoryCheck: true}}, ProductSnapshot: product(), CompanySnapshot: domain.Company{ID: "company", OwnerID: "client", VerificationStatus: "VERIFIED", Data: domain.CompanyData{MonthlyRevenueMinor: 300000000}}, Documents: []domain.Document{{ID: "doc", Type: "BANK_STATEMENT", Status: "UPLOADED"}}}
	r := &memoryRepo{apps: map[string]domain.Aggregate{"app": a}, receipts: map[string]receipt{}}

	return &Service{Repo: r, Companies: companyFake{}}, r
}

func ctx(user, role, key, version string) context.Context {
	return platform.WithHeaders(platform.WithIdentity(context.Background(), platform.Claims{Subject: user, Role: role}), key, version)
}

func statusCode(e error) int {
	var v *platform.Error
	if errors.As(e, &v) {
		return v.Status
	}

	return 0
}

func TestOwnershipAndEmployeePermissions(t *testing.T) {
	s, _ := setup("MANUAL_CHECK")
	if _, e := s.Get(ctx("other", "CLIENT", "key-0001", "1"), "app"); statusCode(e) != 404 {
		t.Fatal("foreign application exposed")
	}

	if _, e := s.Decide(ctx("client", "CLIENT", "key-0001", "1"), "app", domain.ManualDecision{}); statusCode(e) != 403 {
		t.Fatal("client can decide")
	}
}

func TestSubmitAndRequiredDocuments(t *testing.T) {
	s, r := setup("IN_PROGRESS")
	a := r.apps["app"]
	a.Documents = nil
	r.apps["app"] = a
	if _, e := s.Submit(ctx("client", "CLIENT", "submit-001", "1"), "app"); statusCode(e) != 422 {
		t.Fatal("incomplete submit accepted")
	}

	if r.apps["app"].Status != "IN_PROGRESS" {
		t.Fatal("failed validation changed status")
	}

	a.Documents = []domain.Document{{Type: "BANK_STATEMENT", Status: "UPLOADED"}}
	r.apps["app"] = a
	v, e := s.Submit(ctx("client", "CLIENT", "submit-002", "1"), "app")
	if e != nil || v.Status != "SCORING" {
		t.Fatalf("submit failed: %v", e)
	}
}

func TestConcurrentCancelUsesVersion(t *testing.T) {
	s, _ := setup("SCORING")
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, key := range []string{"cancel-001", "cancel-002"} {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			_, e := s.Cancel(ctx("client", "CLIENT", key, "1"), "app")
			results <- e
		}(key)
	}

	wg.Wait()
	close(results)
	success, conflict := 0, 0
	for e := range results {
		if e == nil {
			success++
		} else if statusCode(e) == 409 {
			conflict++
		}
	}

	if success != 1 || conflict != 1 {
		t.Fatalf("success=%d conflict=%d", success, conflict)
	}
}

func TestIssuanceIdempotency(t *testing.T) {
	s, r := setup("ACCEPTED")
	a := r.apps["app"]
	a.AcceptedOfferID = "offer"
	a.Decision = &domain.Decision{Offers: []domain.Offer{{ID: "offer", AmountMinor: 100000000}}}
	r.apps["app"] = a
	first, e := s.Issue(ctx("client", "CLIENT", "issue-001", "1"), "app")
	if e != nil {
		t.Fatal(e)
	}

	second, e := s.Issue(ctx("client", "CLIENT", "issue-001", "1"), "app")
	if e != nil || first.ID != second.ID {
		t.Fatal("repeat created another loan")
	}

	if _, e = s.Issue(ctx("client", "CLIENT", "issue-002", "2"), "app"); statusCode(e) != 409 {
		t.Fatal("second issuance accepted")
	}
}

func TestManualInformationCycle(t *testing.T) {
	s, r := setup("MANUAL_CHECK")
	v, e := s.Decide(ctx("employee", "EMPLOYEE", "decide-001", "1"), "app", domain.ManualDecision{Outcome: "NEEDS_INFO", Reason: "Нужна выписка", Items: []domain.InformationItem{{Title: "Выписка", Type: "BANK_STATEMENT", Required: true}}})
	if e != nil || v.Status != "NEEDS_INFO" {
		t.Fatal(e)
	}

	if Validate(r.apps["app"]).Valid {
		t.Fatal("pending request not required")
	}
}

func TestExpiredOfferAndAlternativeValidation(t *testing.T) {
	s, r := setup("APPROVED")
	a := r.apps["app"]
	a.Decision = &domain.Decision{Offers: []domain.Offer{{ID: "expired", ProductID: a.ProductID, AmountMinor: 100000000, TermMonths: 12, ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339)}}}
	r.apps["app"] = a
	if _, e := s.Accept(ctx("client", "CLIENT", "accept-001", "1"), "app", "expired"); statusCode(e) != 409 {
		t.Fatal("expired offer accepted")
	}
}

func TestScoringDoesNotDirectlyGuaranteeApproval(t *testing.T) {
	s, r := setup("SCORING")
	_ = s
	a := r.apps["app"]
	Finish(&a, domain.Scoring{Score: 800, RiskLevel: "LOW", Affordability: &domain.Affordability{Status: "CALCULATED", MaxAmountMinor: 60000000, Reason: "Платёж не должен превышать 5 000 ₽ в месяц"}})
	if a.Status != "ALTERNATIVE" || a.Decision.Offers[0].AmountMinor != 60000000 || !strings.Contains(a.Decision.Reason, "5 000 ₽") {
		t.Fatal("capacity ignored")
	}

	a = r.apps["app"]
	Finish(&a, domain.Scoring{Score: 800, RiskLevel: "LOW", Affordability: &domain.Affordability{Status: "CALCULATED", MaxAmountMinor: 0}})
	if a.Status != "REJECTED" {
		t.Fatalf("unaffordable loan not rejected: %s", a.Status)
	}

	a = r.apps["app"]
	Finish(&a, domain.Scoring{Score: 800, RiskLevel: "LOW"})
	if a.Status != "MANUAL_CHECK" || a.Decision != nil {
		t.Fatal("decision made without affordability")
	}

	a = r.apps["app"]
	Finish(&a, domain.Scoring{Score: 800, RiskLevel: "LOW", Affordability: &domain.Affordability{Status: "CALCULATED", MaxAmountMinor: 500000000}})
	if a.Status != "APPROVED" || a.Decision.Offers[0].AmountMinor != a.Need.AmountMinor {
		t.Fatalf("affordable loan not approved: %s", a.Status)
	}

	a = r.apps["app"]
	Finish(&a, domain.Scoring{Score: 500, RiskLevel: "MEDIUM"})
	if a.Status != "MANUAL_CHECK" || a.Decision != nil {
		t.Fatal("risk directly became final decision")
	}
}

type memoryStorage struct{}

func (memoryStorage) Put(context.Context, string, io.Reader, string) (int64, error) { return 0, nil }

func (memoryStorage) Get(context.Context, string) (io.ReadCloser, error) { return nil, nil }

func (memoryStorage) Delete(context.Context, string) error { return nil }
