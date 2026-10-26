package usecase

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"robust/application-service/internal/domain"
	platform "robust/platform/core"
)

type Repository interface {
	Get(context.Context, string) (domain.Aggregate, error)
	List(context.Context, string, domain.PageInput) (domain.Page, error)
	Apply(context.Context, platform.Command, string, func(*domain.Aggregate) (domain.Result, error)) (domain.Result, error)
	Claim(context.Context) (*domain.Job, error)
	Complete(context.Context, domain.Job, func(*domain.Aggregate) error, bool) error
	StageObject(context.Context, string) error
	Garbage(context.Context) ([]string, error)
	RemoveObject(context.Context, string) error
}

type CompanyProvider interface {
	Company(context.Context, string) (domain.Company, error)
	Product(context.Context, string) (domain.Product, error)
}

type Scorer interface {
	Evaluate(context.Context, domain.ScoringInput) (domain.Scoring, error)
}

type ObjectStorage interface {
	Put(context.Context, string, io.Reader, string) (int64, error)
	Get(context.Context, string) (io.ReadCloser, error)
	Delete(context.Context, string) error
}

type Service struct {
	Repo      Repository
	Companies CompanyProvider
	Scorer    Scorer
	Storage   ObjectStorage
}

func (s *Service) Get(ctx context.Context, id string) (domain.Aggregate, error) {
	a, e := s.Repo.Get(ctx, id)
	if e != nil {
		return a, e
	}

	u := platform.Identity(ctx)
	if u.Role == "CLIENT" && a.OwnerID != u.Subject {
		return a, platform.NotFound()
	}

	return a, nil
}

func (s *Service) List(ctx context.Context, in domain.PageInput, bank bool) (domain.Page, error) {
	user := platform.Identity(ctx).Subject
	if bank {
		if e := platform.Require(ctx, "EMPLOYEE"); e != nil {
			return domain.Page{}, e
		}

		user = ""
	} else if e := platform.Require(ctx, "CLIENT"); e != nil {
		return domain.Page{}, e
	}

	if in.Offset < 0 || in.Limit < 0 || in.Limit > 100 {
		return domain.Page{}, platform.Invalid("Некорректная пагинация")
	}

	if in.Limit == 0 {
		in.Limit = 20
	}

	return s.Repo.List(ctx, user, in)
}

func event(a *domain.Aggregate, status, message, role string) {
	a.Events = append(a.Events, domain.Event{ID: platform.ID(), At: platform.Now(), FromStatus: a.Status, ToStatus: status, Message: message, ActorRole: role})
	a.Status = status
	a.NextAction = domain.Next(status)
}

func (s *Service) Create(ctx context.Context, in domain.DraftInput) (domain.Application, error) {
	if e := platform.Require(ctx, "CLIENT"); e != nil {
		return domain.Application{}, e
	}

	cmd, e := platform.GetCommand(ctx, "application.create", in, false)
	if e != nil {
		return domain.Application{}, e
	}

	c, e := s.Companies.Company(ctx, in.Need.CompanyID)
	if e != nil {
		return domain.Application{}, e
	}

	if c.OwnerID != cmd.User {
		return domain.Application{}, platform.NotFound()
	}

	p, e := s.Companies.Product(ctx, in.ProductID)
	if e != nil {
		return domain.Application{}, e
	}

	if e = validNeed(in.Need, p); e != nil {
		return domain.Application{}, e
	}

	r, e := s.Repo.Apply(ctx, cmd, "", func(a *domain.Aggregate) (domain.Result, error) {
		id := platform.ID()
		a.Application = domain.Application{ID: id, CompanyID: c.ID, ProductID: p.ID, Need: in.Need, Consents: in.Consents, OwnerID: cmd.User, CompanyName: c.Data.Name, INN: c.Data.INN, ProductName: p.Name, Number: "K-" + time.Now().Format("2006") + "-" + strings.ToUpper(id[:8]), MissingFields: []string{}}
		a.CompanySnapshot = c
		a.ProductSnapshot = p
		a.Documents = []domain.Document{}
		a.InformationRequests = []domain.InformationRequest{}
		event(a, "IN_PROGRESS", "Черновик заявки создан", cmd.Role)
		return domain.Result{}, nil
	})
	if e != nil {
		return domain.Application{}, e
	}

	return *r.Application, nil
}

func validNeed(n domain.Need, p domain.Product) error {
	if n.CompanyID == "" || n.AmountMinor < p.MinAmountMinor || n.AmountMinor > p.MaxAmountMinor || n.TermMonths < p.MinTermMonths || n.TermMonths > p.MaxTermMonths || !domain.Supports(p, n.Purpose) || len(n.Description) > 2000 {
		return platform.Invalid("Параметры заявки не соответствуют выбранному продукту")
	}

	return nil
}

func (s *Service) mutate(ctx context.Context, id, op string, payload any, role string, fn func(*domain.Aggregate) (domain.Result, error)) (domain.Result, error) {
	if e := platform.Require(ctx, role); e != nil {
		return domain.Result{}, e
	}

	cmd, e := platform.GetCommand(ctx, op+"/"+id, payload, true)
	if e != nil {
		return domain.Result{}, e
	}

	return s.Repo.Apply(ctx, cmd, id, func(a *domain.Aggregate) (domain.Result, error) {
		if role == "CLIENT" && a.OwnerID != cmd.User {
			return domain.Result{}, platform.NotFound()
		}

		if e := platform.CheckVersion(cmd, a.Version); e != nil {
			return domain.Result{}, e
		}

		return fn(a)
	})
}

func (s *Service) Patch(ctx context.Context, id string, in domain.Patch) (domain.Application, error) {
	var c *domain.Company
	var p *domain.Product
	if in.Need != nil {
		v, e := s.Companies.Company(ctx, in.Need.CompanyID)
		if e != nil {
			return domain.Application{}, e
		}

		if v.OwnerID != platform.Identity(ctx).Subject {
			return domain.Application{}, platform.NotFound()
		}

		c = &v
	}

	if in.ProductID != nil {
		v, e := s.Companies.Product(ctx, *in.ProductID)
		if e != nil {
			return domain.Application{}, e
		}

		p = &v
	}

	r, e := s.mutate(ctx, id, "application.patch", in, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if !domain.CanEdit(a.Status) {
			return domain.Result{}, platform.Conflict("Изменение заявки недоступно")
		}

		if c != nil && c.ID != a.CompanyID {
			return domain.Result{}, platform.Invalid("Компания заявки неизменяема; создайте новый черновик")
		}

		if c != nil {
			a.CompanySnapshot = *c
		}

		if p != nil {
			if a.Submitted && p.ID != a.ProductID {
				return domain.Result{}, platform.Invalid("После отправки продукт неизменяем")
			}

			a.ProductSnapshot = *p
			a.ProductID = p.ID
			a.ProductName = p.Name
		}

		if in.Need != nil {
			a.Need = *in.Need
		}

		if in.Consents != nil {
			a.Consents = in.Consents
		}

		if e := validNeed(a.Need, a.ProductSnapshot); e != nil {
			return domain.Result{}, e
		}

		event(a, a.Status, "Данные заявки обновлены", "CLIENT")
		return domain.Result{}, nil
	})
	if e != nil {
		return domain.Application{}, e
	}

	return *r.Application, nil
}

func Validate(a domain.Aggregate) domain.Validation {
	v := domain.Validation{MissingFields: []string{}, Errors: []domain.FieldError{}}
	missing := func(field, msg string) {
		v.MissingFields = append(v.MissingFields, field)
		v.Errors = append(v.Errors, domain.FieldError{Field: field, Message: msg})
	}
	if a.Consents == nil || !a.Consents.PersonalDataProcessing || !a.Consents.CreditHistoryCheck {
		missing("consents", "Подтвердите согласия")
	}

	if a.CompanySnapshot.VerificationStatus != "VERIFIED" {
		missing("company.verification", "Подтвердите сведения компании")
	}

	for _, typ := range a.ProductSnapshot.RequiredDocumentTypes {
		found := false
		for _, d := range a.Documents {
			if d.Type == typ && d.Status == "UPLOADED" {
				found = true
			}
		}

		if !found {
			missing("documents."+typ, "Загрузите обязательный документ")
		}
	}

	for _, req := range a.InformationRequests {
		for _, item := range req.Items {
			if item.Required && item.Status != "FULFILLED" {
				missing("informationRequests."+item.ID, item.Title)
			}
		}
	}

	v.Valid = len(v.Errors) == 0

	return v
}

func (s *Service) Validation(ctx context.Context, id string) (domain.Validation, error) {
	if _, e := platform.GetCommand(ctx, "application.validation/"+id, id, true); e != nil {
		return domain.Validation{}, e
	}

	a, e := s.Get(ctx, id)
	if e != nil {
		return domain.Validation{}, e
	}

	if e := platform.Require(ctx, "CLIENT"); e != nil {
		return domain.Validation{}, e
	}

	cmd, _ := platform.GetCommand(ctx, "application.validation/"+id, id, true)
	if e = platform.CheckVersion(cmd, a.Version); e != nil {
		return domain.Validation{}, e
	}

	return Validate(a), nil
}

func (s *Service) Submit(ctx context.Context, id string) (domain.Application, error) {
	old, e := s.Get(ctx, id)
	if e != nil {
		return domain.Application{}, e
	}

	c, e := s.Companies.Company(ctx, old.CompanyID)
	if e != nil {
		return domain.Application{}, e
	}

	p, e := s.Companies.Product(ctx, old.ProductID)
	if e != nil {
		return domain.Application{}, e
	}

	r, e := s.mutate(ctx, id, "application.submit", id, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if !domain.CanEdit(a.Status) {
			return domain.Result{}, platform.Conflict("Отправка недоступна в текущем статусе")
		}

		a.CompanySnapshot = c
		a.ProductSnapshot = p
		a.CompanyName = c.Data.Name
		if e := validNeed(a.Need, p); e != nil {
			return domain.Result{}, e
		}

		v := Validate(*a)
		if !v.Valid {
			return domain.Result{}, platform.Invalid(strings.Join(v.MissingFields, ", "))
		}

		a.Decision = nil
		a.Scoring = nil
		a.Recommendation = ""
		a.MissingFields = []string{}
		a.Submitted = true
		event(a, "SCORING", "Заявка отправлена на проверку", "CLIENT")
		return domain.Result{}, nil
	})
	if e != nil {
		return domain.Application{}, e
	}

	return *r.Application, nil
}

func (s *Service) Cancel(ctx context.Context, id string) (domain.Application, error) {
	r, e := s.mutate(ctx, id, "application.cancel", id, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if !domain.CanEdit(a.Status) && a.Status != "SCORING" && a.Status != "MANUAL_CHECK" {
			return domain.Result{}, platform.Conflict("Отзыв недоступен")
		}

		event(a, "CANCELLED", "Заявка отозвана клиентом", "CLIENT")
		return domain.Result{}, nil
	})
	if e != nil {
		return domain.Application{}, e
	}

	return *r.Application, nil
}

func (s *Service) Accept(ctx context.Context, id, offerID string) (domain.Application, error) {
	r, e := s.mutate(ctx, id, "application.accept", offerID, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if (a.Status != "APPROVED" && a.Status != "ALTERNATIVE") || a.Decision == nil {
			return domain.Result{}, platform.Conflict("Нет доступного предложения")
		}

		for _, o := range a.Decision.Offers {
			if o.ID == offerID && domain.OfferValid(o, a.ProductSnapshot) {
				a.AcceptedOfferID = offerID
				event(a, "ACCEPTED", "Клиент принял предложение", "CLIENT")
				return domain.Result{}, nil
			}
		}

		return domain.Result{}, platform.Conflict("Предложение недействительно или истекло")
	})
	if e != nil {
		return domain.Application{}, e
	}

	return *r.Application, nil
}

func (s *Service) Decline(ctx context.Context, id, reason string) (domain.Application, error) {
	r, e := s.mutate(ctx, id, "application.decline", reason, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if a.Status != "APPROVED" && a.Status != "ALTERNATIVE" {
			return domain.Result{}, platform.Conflict("Нет предложения для отказа")
		}

		if len(reason) > 1000 {
			return domain.Result{}, platform.Invalid("Комментарий длиннее 1000 символов")
		}

		message := "Клиент отказался от условий"
		if reason != "" {
			message += ": " + reason
		}

		event(a, "DECLINED", message, "CLIENT")
		return domain.Result{}, nil
	})
	if e != nil {
		return domain.Application{}, e
	}

	return *r.Application, nil
}

func (s *Service) Issue(ctx context.Context, id string) (domain.Loan, error) {
	r, e := s.mutate(ctx, id, "application.issue", id, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if a.Status != "ACCEPTED" || a.Loan != nil || a.Decision == nil {
			return domain.Result{}, platform.Conflict("Выдача недоступна")
		}

		for _, o := range a.Decision.Offers {
			if o.ID == a.AcceptedOfferID {
				a.Loan = &domain.Loan{ID: platform.ID(), ApplicationID: id, OfferID: o.ID, AmountMinor: o.AmountMinor, IssuedAt: platform.Now(), Simulated: true}
				event(a, "ISSUED", "Подписание и выдача смоделированы; реальные деньги не перечисляются", "CLIENT")
				return domain.Result{Loan: a.Loan}, nil
			}
		}

		return domain.Result{}, platform.Conflict("Предложение не найдено")
	})
	if e != nil {
		return domain.Loan{}, e
	}

	return *r.Loan, nil
}

func (s *Service) Decide(ctx context.Context, id string, in domain.ManualDecision) (domain.Application, error) {
	r, e := s.mutate(ctx, id, "application.decide", in, "EMPLOYEE", func(a *domain.Aggregate) (domain.Result, error) {
		if a.Status != "MANUAL_CHECK" {
			return domain.Result{}, platform.Conflict("Заявка не находится на ручной проверке")
		}

		if len(strings.TrimSpace(in.Reason)) == 0 || len(in.Reason) > 2000 {
			return domain.Result{}, platform.Invalid("Укажите причину решения")
		}

		if in.Outcome == "NEEDS_INFO" {
			items := in.Items
			for _, field := range in.RequestedFields {
				items = append(items, domain.InformationItem{Title: field, Type: "OTHER", Required: true})
			}

			if len(items) == 0 {
				return domain.Result{}, platform.Invalid("Укажите требуемые сведения")
			}

			if len(in.Offers) > 0 {
				return domain.Result{}, platform.Invalid("Запрос данных не содержит предложений")
			}

			a.MissingFields = []string{}
			for i := range items {
				if strings.TrimSpace(items[i].Title) == "" {
					return domain.Result{}, platform.Invalid("Нужно название запрашиваемого документа")
				}

				if !validDocumentType(items[i].Type) {
					return domain.Result{}, platform.Invalid("Неизвестный тип документа")
				}

				items[i].ID = platform.ID()
				items[i].Status = "PENDING"
				if items[i].Required {
					a.MissingFields = append(a.MissingFields, items[i].ID)
				}
			}

			a.InformationRequests = append(a.InformationRequests, domain.InformationRequest{ID: platform.ID(), Message: in.Reason, Items: items, CreatedAt: platform.Now()})
			event(a, "NEEDS_INFO", "Банк запросил дополнительные сведения", "EMPLOYEE")
			return domain.Result{}, nil
		}

		if in.Outcome != "APPROVED" && in.Outcome != "ALTERNATIVE" && in.Outcome != "REJECTED" {
			return domain.Result{}, platform.Invalid("Неизвестное решение")
		}

		if in.Outcome == "REJECTED" && len(in.Offers) > 0 {
			return domain.Result{}, platform.Invalid("Отказ не содержит предложений")
		}

		if in.Outcome != "REJECTED" && len(in.Offers) == 0 {
			return domain.Result{}, platform.Invalid("Требуется предложение")
		}

		for i := range in.Offers {
			in.Offers[i].ID = platform.ID()
			if in.Offers[i].ProductID == "" {
				in.Offers[i].ProductID = a.ProductID
			}

			if in.Offers[i].ExpiresAt == "" {
				in.Offers[i].ExpiresAt = time.Now().UTC().Add(7 * 24 * time.Hour).Format(time.RFC3339Nano)
			}

			if !domain.OfferValid(in.Offers[i], a.ProductSnapshot) {
				return domain.Result{}, platform.Invalid("Предложение выходит за условия продукта или истекло")
			}

			if in.Outcome == "APPROVED" && (in.Offers[i].AmountMinor != a.Need.AmountMinor || in.Offers[i].TermMonths != a.Need.TermMonths) {
				return domain.Result{}, platform.Invalid("Для других суммы или срока используйте ALTERNATIVE")
			}

			in.Offers[i].MonthlyPaymentMinor = domain.Payment(in.Offers[i].AmountMinor, in.Offers[i].TermMonths, in.Offers[i].AnnualRateBps)
		}

		a.Decision = &domain.Decision{ID: platform.ID(), Outcome: in.Outcome, Reason: in.Reason, Offers: in.Offers, DecidedAt: platform.Now()}
		event(a, "CHECKED", "Ручная проверка завершена", "EMPLOYEE")
		event(a, in.Outcome, "Банк сформировал решение", "EMPLOYEE")
		return domain.Result{}, nil
	})
	if e != nil {
		return domain.Application{}, e
	}

	return *r.Application, nil
}

func validDocumentType(v string) bool {
	switch v {
	case "FINANCIAL_REPORT", "BANK_STATEMENT", "REGISTRATION", "OTHER":
		return true
	}

	return false
}

func (s *Service) Upload(ctx context.Context, id string, d domain.Document, reader io.Reader) (domain.Document, error) {
	if e := platform.Require(ctx, "CLIENT"); e != nil {
		return d, e
	}

	_, e := s.Get(ctx, id)
	if e != nil {
		return d, e
	}

	if !validDocumentType(d.Type) {
		return d, platform.Invalid("Неизвестный тип документа")
	}

	if _, e = platform.GetCommand(ctx, "document.upload/"+id, d, true); e != nil {
		return d, e
	}

	b := bufio.NewReader(reader)
	head, e := b.Peek(512)
	if e != nil && e != io.EOF && e != bufio.ErrBufferFull {
		return d, e
	}

	actual := http.DetectContentType(head)
	if actual != "application/pdf" && actual != "image/png" && actual != "image/jpeg" {
		return d, platform.E("UNSUPPORTED_MEDIA_TYPE", "Разрешены PDF, PNG и JPEG", 415)
	}

	if actual != d.ContentType {
		return d, platform.E("UNSUPPORTED_MEDIA_TYPE", "Содержимое не соответствует MIME", 415)
	}

	d.ID = platform.ID()
	d.ApplicationID = id
	d.FileName = filepath.Base(strings.ReplaceAll(d.FileName, "\\", "/"))
	if d.FileName == "." || d.FileName == "" || len(d.FileName) > 255 {
		return d, platform.Invalid("Некорректное имя файла")
	}

	d.ObjectKey = "documents/" + d.ID
	d.Status = "UPLOADED"
	if e = s.Repo.StageObject(ctx, d.ObjectKey); e != nil {
		return d, e
	}

	hash := sha256.New()
	d.SizeBytes, e = s.Storage.Put(ctx, d.ObjectKey, io.TeeReader(io.LimitReader(b, 10485761), hash), d.ContentType)
	if e != nil {
		return d, e
	}

	if d.SizeBytes > 10485760 {
		return d, platform.E("FILE_TOO_LARGE", "Файл больше 10 МиБ", 413)
	}

	d.Checksum = hex.EncodeToString(hash.Sum(nil))
	if d.SizeBytes == 0 {
		return d, platform.Invalid("Пустой файл")
	}

	payload := d
	payload.ID = ""
	payload.ObjectKey = ""
	r, e := s.mutate(ctx, id, "document.upload", payload, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if !domain.CanEdit(a.Status) {
			return domain.Result{}, platform.Conflict("Загрузка недоступна")
		}

		if d.RequestItemID != "" {
			found := false
			for i := range a.InformationRequests {
				for j := range a.InformationRequests[i].Items {
					item := &a.InformationRequests[i].Items[j]
					if item.ID == d.RequestItemID {
						if item.Type != d.Type {
							return domain.Result{}, platform.Invalid("Тип документа не соответствует запросу")
						}

						item.Status = "FULFILLED"
						found = true
					}
				}
			}

			if !found {
				return domain.Result{}, platform.NotFound()
			}
		}

		if len(a.Documents) >= 100 {
			return domain.Result{}, platform.Invalid("Не больше 100 документов на заявку")
		}

		a.Documents = append(a.Documents, d)
		event(a, a.Status, "Документ прикреплён", "CLIENT")
		return domain.Result{Document: &d}, nil
	})
	if e != nil {
		return d, e
	}

	if r.Document.ID != d.ID {
		_ = s.Storage.Delete(ctx, d.ObjectKey)
		_ = s.Repo.RemoveObject(ctx, d.ObjectKey)
	}

	return *r.Document, nil
}

func (s *Service) DeleteDocument(ctx context.Context, id, documentID string) error {
	_, e := s.mutate(ctx, id, "document.delete", documentID, "CLIENT", func(a *domain.Aggregate) (domain.Result, error) {
		if !domain.CanEdit(a.Status) {
			return domain.Result{}, platform.Conflict("Удаление недоступно")
		}

		for i, d := range a.Documents {
			if d.ID == documentID {
				a.Documents = append(a.Documents[:i], a.Documents[i+1:]...)
				if d.RequestItemID != "" {
					for x := range a.InformationRequests {
						for y := range a.InformationRequests[x].Items {
							item := &a.InformationRequests[x].Items[y]
							if item.ID == d.RequestItemID {
								item.Status = "PENDING"
								for _, other := range a.Documents {
									if other.RequestItemID == item.ID {
										item.Status = "FULFILLED"
									}
								}
							}
						}
					}
				}

				event(a, a.Status, "Документ удалён", "CLIENT")
				return domain.Result{}, nil
			}
		}

		return domain.Result{}, platform.NotFound()
	})

	return e
}

func (s *Service) Download(ctx context.Context, id, docID string) (domain.Document, io.ReadCloser, error) {
	a, e := s.Get(ctx, id)
	if e != nil {
		return domain.Document{}, nil, e
	}

	for _, d := range a.Documents {
		if d.ID == docID {
			r, e := s.Storage.Get(ctx, d.ObjectKey)

			return d, r, e
		}
	}

	return domain.Document{}, nil, platform.NotFound()
}

func (s *Service) Process(ctx context.Context) error {
	job, e := s.Repo.Claim(ctx)
	if e != nil || job == nil {
		return e
	}

	a, e := s.Repo.Get(ctx, job.ApplicationID)
	if e != nil {
		return e
	}

	f := a.CompanySnapshot.Data
	score, scoreErr := s.Scorer.Evaluate(ctx, domain.ScoringInput{ScoringRequestID: job.ID, ApplicationID: a.ID, ApplicationVersion: job.Revision, Features: domain.Features{BusinessAgeMonths: f.BusinessAgeMonths, MonthlyRevenueMinor: f.MonthlyRevenueMinor, MonthlyExpensesMinor: f.MonthlyExpensesMinor, MonthlyDebtPaymentsMinor: f.MonthlyDebtPaymentsMinor, RequestedAmountMinor: a.Need.AmountMinor, TermMonths: a.Need.TermMonths, HasOverdueDebt: f.HasOverdueDebt}})
	if scoreErr != nil && job.Attempts < 3 {
		return s.Repo.Complete(ctx, *job, nil, true)
	}

	return s.Repo.Complete(ctx, *job, func(current *domain.Aggregate) error {
		if scoreErr != nil {
			TechnicalReview(current)
		} else {
			Finish(current, score)
		}

		return nil
	}, false)
}

func (s *Service) Reconcile(ctx context.Context) error {
	keys, e := s.Repo.Garbage(ctx)
	if e != nil {
		return e
	}

	for _, key := range keys {
		if e = s.Storage.Delete(ctx, key); e != nil {
			return e
		}

		if e = s.Repo.RemoveObject(ctx, key); e != nil {
			return e
		}
	}

	return nil
}

func Finish(a *domain.Aggregate, score domain.Scoring) {
	a.Scoring = &score
	a.Scoring.Checks = []domain.Check{{Source: "FNS_MOCK", Status: "PASSED", Summary: "Локальная проверка компании", CheckedAt: platform.Now()}, {Source: "BKI_MOCK", Status: "PASSED", Summary: "Синтетическая кредитная история", CheckedAt: platform.Now()}}
	a.Recommendation = score.Recommendation
	if score.RiskLevel != "LOW" {
		event(a, "MANUAL_CHECK", "Требуется решение сотрудника", "SYSTEM")
		return
	}

	amount := a.Need.AmountMinor
	capacity := a.CompanySnapshot.Data.MonthlyRevenueMinor * 6
	outcome := "APPROVED"
	reason := "Риск и условия продукта соответствуют правилам демонстрационной модели"
	offers := []domain.Offer{}
	if capacity < amount {
		amount = capacity
		outcome = "ALTERNATIVE"
		reason = "Предложена меньшая сумма с учётом оборота бизнеса"
	}

	if amount < a.ProductSnapshot.MinAmountMinor {
		outcome = "REJECTED"
		reason = "Доступная сумма ниже минимального размера продукта"
	} else {
		offers = append(offers, domain.Offer{ID: platform.ID(), ProductID: a.ProductID, AmountMinor: amount, TermMonths: a.Need.TermMonths, AnnualRateBps: a.ProductSnapshot.AnnualRateBps, MonthlyPaymentMinor: domain.Payment(amount, a.Need.TermMonths, a.ProductSnapshot.AnnualRateBps), ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour).Format(time.RFC3339Nano), Collateral: a.ProductSnapshot.Collateral})
	}

	a.Decision = &domain.Decision{ID: platform.ID(), Outcome: outcome, Reason: reason, Offers: offers, DecidedAt: platform.Now()}
	event(a, "CHECKED", "Проверки завершены", "SYSTEM")
	event(a, outcome, "Сформировано решение банка", "SYSTEM")
}

func TechnicalReview(a *domain.Aggregate) {
	a.Recommendation = "MANUAL_REVIEW"
	event(a, "MANUAL_CHECK", "Автоматическая проверка недоступна; требуется сотрудник", "SYSTEM")
}
