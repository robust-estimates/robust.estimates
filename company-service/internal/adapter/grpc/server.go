package grpcadapter

import (
	"context"

	"robust/company-service/internal/domain"
	"robust/company-service/internal/usecase"
	pb "robust/contracts/credit/v1"
	"robust/platform"
)

type Server struct {
	pb.UnimplementedCompanyServiceServer
	Service *usecase.Service
}

func (s *Server) CreateCompany(ctx context.Context, r *pb.CompanyInput) (*pb.Company, error) {
	var in domain.Input
	if e := platform.FromProto(r, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.Create(ctx, in)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Company{}), nil
}

func (s *Server) GetCompany(ctx context.Context, r *pb.ID) (*pb.Company, error) {
	v, e := s.Service.Get(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Company{}), nil
}

func (s *Server) ListCompanies(ctx context.Context, r *pb.Empty) (*pb.CompanyList, error) {
	v, e := s.Service.List(ctx)
	if e != nil {
		return nil, e
	}

	return platform.PB(map[string]any{"items": v}, &pb.CompanyList{}), nil
}

func (s *Server) UpdateCompany(ctx context.Context, r *pb.CompanyUpdate) (*pb.Company, error) {
	var in domain.Input
	if r.Data == nil {
		return nil, platform.Invalid("Нужны данные компании")
	}

	if e := platform.FromProto(r.Data, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.Update(ctx, r.Id, in)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Company{}), nil
}

func (s *Server) VerifyCompany(ctx context.Context, r *pb.ID) (*pb.Company, error) {
	v, e := s.Service.Verify(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Company{}), nil
}

func (s *Server) ListProducts(ctx context.Context, r *pb.Empty) (*pb.ProductList, error) {
	v, e := s.Service.Products(ctx)
	if e != nil {
		return nil, e
	}

	return platform.PB(map[string]any{"items": v}, &pb.ProductList{}), nil
}

func (s *Server) GetProduct(ctx context.Context, r *pb.ID) (*pb.Product, error) {
	v, e := s.Service.Product(ctx, r.Id)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Product{}), nil
}

func (s *Server) Recommend(ctx context.Context, r *pb.RecommendationInput) (*pb.RecommendationList, error) {
	var in domain.RecommendationInput
	if e := platform.FromProto(r, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.Recommend(ctx, in)
	if e != nil {
		return nil, e
	}

	return platform.PB(map[string]any{"items": v}, &pb.RecommendationList{}), nil
}

func (s *Server) Calculate(ctx context.Context, r *pb.CalculationInput) (*pb.Calculation, error) {
	var in domain.CalculationInput
	if e := platform.FromProto(r, &in); e != nil {
		return nil, e
	}

	v, e := s.Service.Calculate(ctx, in)
	if e != nil {
		return nil, e
	}

	return platform.PB(v, &pb.Calculation{}), nil
}
