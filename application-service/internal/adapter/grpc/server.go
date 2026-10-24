package grpcadapter

import (
	"context"
	"io"

	"google.golang.org/grpc"

	"robust/application-service/internal/domain"
	"robust/application-service/internal/usecase"
	pb "robust/contracts/credit/v1"
	"robust/platform"
)

type Server struct {
	pb.UnimplementedApplicationServiceServer
	Service *usecase.Service
}

func (s *Server) CreateApplication(ctx context.Context, r *pb.DraftInput) (*pb.Application, error) {
	var in domain.DraftInput
	if e := platform.FromProto(r, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.Create(ctx, in)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Application{}), nil
}

func (s *Server) GetApplication(ctx context.Context, r *pb.ID) (*pb.Application, error) {
	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v.Application, &pb.Application{}), nil
}

func (s *Server) ListApplications(ctx context.Context, r *pb.PageInput) (*pb.ApplicationPage, error) {
	return s.list(ctx, r, false)
}

func (s *Server) BankApplications(ctx context.Context, r *pb.PageInput) (*pb.ApplicationPage, error) {
	return s.list(ctx, r, true)
}

func (s *Server) list(ctx context.Context, r *pb.PageInput, bank bool) (*pb.ApplicationPage, error) {
	var in domain.PageInput
	if e := platform.FromProto(r, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.List(ctx, in, bank)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.ApplicationPage{}), nil
}

func (s *Server) PatchApplication(ctx context.Context, r *pb.PatchInput) (*pb.Application, error) {
	if r.Patch == nil {
		return nil, platform.Invalid("Пустое изменение")
	}

	var in domain.Patch
	if e := platform.FromProto(r.Patch, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.Patch(ctx, r.Id, in)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Application{}), nil
}

func (s *Server) Validate(ctx context.Context, r *pb.ID) (*pb.Validation, error) {
	v, e := s.Service.Validation(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Validation{}), nil
}

func (s *Server) Submit(ctx context.Context, r *pb.ID) (*pb.Application, error) {
	v, e := s.Service.Submit(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Application{}), nil
}

func (s *Server) Cancel(ctx context.Context, r *pb.ID) (*pb.Application, error) {
	v, e := s.Service.Cancel(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Application{}), nil
}

func (s *Server) Accept(ctx context.Context, r *pb.AcceptInput) (*pb.Application, error) {
	v, e := s.Service.Accept(ctx, r.Id, r.OfferId)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Application{}), nil
}

func (s *Server) Decline(ctx context.Context, r *pb.DeclineInput) (*pb.Application, error) {
	v, e := s.Service.Decline(ctx, r.Id, r.Reason)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Application{}), nil
}

func (s *Server) Issue(ctx context.Context, r *pb.ID) (*pb.Loan, error) {
	v, e := s.Service.Issue(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Loan{}), nil
}

func (s *Server) GetDecision(ctx context.Context, r *pb.ID) (*pb.Decision, error) {
	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	if v.Decision == nil {
		return nil, platform.Conflict("Решение ещё не сформировано")
	}

	return platform.PB(v.Decision, &pb.Decision{}), nil
}

func (s *Server) History(ctx context.Context, r *pb.ID) (*pb.EventList, error) {
	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(map[string]any{"items": v.Events}, &pb.EventList{}), nil
}

func (s *Server) GetScoring(ctx context.Context, r *pb.ID) (*pb.Scoring, error) {
	if e := platform.Require(ctx, "EMPLOYEE"); e != nil {
		return nil, e
	}

	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	if v.Scoring == nil {
		return nil, platform.Conflict("Скоринг ещё не готов")
	}

	return platform.PB(v.Scoring, &pb.Scoring{}), nil
}

func (s *Server) Decide(ctx context.Context, r *pb.ManualInput) (*pb.Application, error) {
	if r.Decision == nil {
		return nil, platform.Invalid("Нужно решение")
	}

	var in domain.ManualDecision
	if e := platform.FromProto(r.Decision, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.Decide(ctx, r.Id, in)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Application{}), nil
}

func (s *Server) GetOverview(ctx context.Context, r *pb.ID) (*pb.Overview, error) {
	if e := platform.Require(ctx, "EMPLOYEE"); e != nil {
		return nil, e
	}

	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(map[string]any{"application": v.Application, "companySnapshot": v.CompanySnapshot, "productSnapshot": v.ProductSnapshot, "scoring": v.Scoring, "decision": v.Decision, "documents": v.Documents, "informationRequests": v.InformationRequests}, &pb.Overview{}), nil
}

func (s *Server) InformationRequests(ctx context.Context, r *pb.ID) (*pb.InformationList, error) {
	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(map[string]any{"items": v.InformationRequests}, &pb.InformationList{}), nil
}

func (s *Server) Feed(ctx context.Context, r *pb.Empty) (*pb.FeedList, error) {
	if e := platform.Require(ctx, "CLIENT"); e != nil {
		return nil, e
	}

	v, e := s.Service.List(ctx, domain.PageInput{Limit: 100}, false)
	if e != nil {
		return nil, e
	}

	items := []domain.FeedItem{}
	for _, a := range v.Items {
		items = append(items, domain.FeedItem{ID: a.ID + "-" + a.Status, Title: a.ProductName + ": " + a.Status, Body: a.NextAction, Reason: "Изменение вашей заявки " + a.Number, CreatedAt: a.UpdatedAt, ApplicationID: a.ID})
	}

	items = append(items, domain.FeedItem{ID: "business-news-v1", Title: "Подготовьте финансовую отчётность заранее", Body: "Для рассмотрения кредита подготовьте актуальную отчётность и выписку по счёту.", Reason: "Рекомендация для предпринимателей", CreatedAt: "2026-10-01T00:00:00Z"})

	return platform.PB(map[string]any{"items": items}, &pb.FeedList{}), nil
}

func (s *Server) ListDocuments(ctx context.Context, r *pb.ID) (*pb.DocumentList, error) {
	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(map[string]any{"items": v.Documents}, &pb.DocumentList{}), nil
}

type uploadReader struct {
	stream grpc.ClientStreamingServer[pb.UploadChunk, pb.Document]
	buffer []byte
	total  int64
}

func (r *uploadReader) Read(p []byte) (int, error) {
	for len(r.buffer) == 0 {
		v, e := r.stream.Recv()
		if e != nil {
			return 0, e
		}

		chunk, ok := v.Value.(*pb.UploadChunk_Data)
		if !ok {
			return 0, platform.Invalid("После метаданных ожидаются только данные")
		}

		r.total += int64(len(chunk.Data))
		if r.total > 10485760 {
			return 0, platform.E("FILE_TOO_LARGE", "Файл больше 10 МиБ", 413)
		}

		r.buffer = chunk.Data
	}

	n := copy(p, r.buffer)
	r.buffer = r.buffer[n:]

	return n, nil
}

func (s *Server) UploadDocument(stream grpc.ClientStreamingServer[pb.UploadChunk, pb.Document]) error {
	first, e := stream.Recv()
	if e != nil {
		return e
	}

	m := first.GetMetadata()
	if m == nil {
		return platform.Invalid("Первое сообщение должно содержать метаданные")
	}

	v, e := s.Service.Upload(stream.Context(), m.ApplicationId, domain.Document{Type: m.Type, FileName: m.FileName, ContentType: m.ContentType, RequestItemID: m.RequestItemId}, &uploadReader{stream: stream})
	if e != nil {
		return e
	}

	return stream.SendAndClose(platform.PB(v, &pb.Document{}))
}

func (s *Server) DeleteDocument(ctx context.Context, r *pb.DocumentID) (*pb.Empty, error) {
	e := s.Service.DeleteDocument(ctx, r.ApplicationId, r.DocumentId)

	return &pb.Empty{}, e
}

func (s *Server) DownloadDocument(r *pb.DocumentID, stream grpc.ServerStreamingServer[pb.DownloadChunk]) error {
	d, reader, e := s.Service.Download(stream.Context(), r.ApplicationId, r.DocumentId)
	if e != nil {
		return e
	}

	defer reader.Close()
	if e = stream.Send(&pb.DownloadChunk{FileName: d.FileName, ContentType: d.ContentType}); e != nil {
		return e
	}

	buf := make([]byte, 64*1024)
	for {
		n, e := reader.Read(buf)
		if n > 0 {
			if se := stream.Send(&pb.DownloadChunk{Data: buf[:n]}); se != nil {
				return se
			}
		}

		if e == io.EOF {
			return nil
		}

		if e != nil {
			return e
		}
	}
}
