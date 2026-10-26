package grpcadapter

import (
	"context"
	"time"

	"robust/application-service/internal/domain"
	pb "robust/contracts/credit/v1"
	"robust/platform"
)

type CompanyClient struct{ Client pb.CompanyServiceClient }

func (c *CompanyClient) Company(ctx context.Context, id string) (domain.Company, error) {
	ctx, cancel := context.WithTimeout(platform.Outgoing(ctx), 5*time.Second)
	defer cancel()
	v, e := c.Client.GetCompany(ctx, &pb.ID{Id: id})
	if e != nil {
		return domain.Company{}, e
	}

	var out domain.Company
	e = platform.FromProto(v, &out)

	return out, e
}

func (c *CompanyClient) Product(ctx context.Context, id string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(platform.Outgoing(ctx), 5*time.Second)
	defer cancel()
	v, e := c.Client.GetProduct(ctx, &pb.ID{Id: id})
	if e != nil {
		return domain.Product{}, e
	}

	var out domain.Product
	e = platform.FromProto(v, &out)

	return out, e
}

type ScoringClient struct{ Client pb.ScoringServiceClient }

func (c *ScoringClient) Evaluate(ctx context.Context, in domain.ScoringInput) (domain.Scoring, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	v, e := c.Client.Evaluate(ctx, platform.PB(in, &pb.ScoringInput{}))
	if e != nil {
		return domain.Scoring{}, e
	}

	if v.ApplicationId != in.ApplicationID || v.ApplicationVersion != in.ApplicationVersion || v.ScoringRequestId != in.ScoringRequestID || v.Scoring == nil {
		return domain.Scoring{}, platform.Invalid("Ответ скоринга не соответствует запросу")
	}

	var out domain.Scoring
	e = platform.FromProto(v.Scoring, &out)
	if out.Score < 0 || out.Score > 1000 || (out.RiskLevel != "LOW" && out.RiskLevel != "MEDIUM" && out.RiskLevel != "HIGH") {
		return out, platform.Invalid("Некорректная оценка риска")
	}

	return out, e
}
