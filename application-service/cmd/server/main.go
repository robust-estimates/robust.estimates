package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"

	"robust/application-service/internal/adapter/external"
	grpcadapter "robust/application-service/internal/adapter/grpc"
	"robust/application-service/internal/adapter/postgres"
	"robust/application-service/internal/adapter/s3"
	"robust/application-service/internal/usecase"
	pb "robust/contracts/credit/v1"
	"robust/platform"
)

func main() {
	if platform.Healthcheck(platform.Env("GRPC_ADDR", ":9092")) {
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "init-storage" {
		store, e := s3.New(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), platform.Env("S3_BUCKET", "credit-documents"), os.Getenv("S3_SECURE") == "true")
		if e != nil {
			panic(e)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		exists, e := store.Client.BucketExists(ctx, store.Bucket)
		if e != nil {
			panic(e)
		}

		if !exists {
			if e = store.Client.MakeBucket(ctx, store.Bucket, minio.MakeBucketOptions{}); e != nil {
				panic(e)
			}
		}

		return
	}

	ctx := context.Background()
	db := platform.OpenDB(ctx, os.Getenv("DATABASE_URL"))
	defer db.Close()
	e := platform.Migrate(ctx, db, platform.Env("MIGRATIONS_DIR", "migrations"))
	if e != nil {
		panic(e)
	}

	storage, e := s3.New(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), platform.Env("S3_BUCKET", "credit-documents"), os.Getenv("S3_SECURE") == "true")
	if e != nil {
		panic(e)
	}

	companyConn := platform.Dial(platform.Env("COMPANY_GRPC_ADDR", "localhost:9091"))
	defer companyConn.Close()
	scoringConn := platform.Dial(platform.Env("SCORING_GRPC_ADDR", "localhost:9093"))
	defer scoringConn.Close()
	secret := platform.Secret()
	service := &usecase.Service{Repo: &postgres.Repository{DB: db}, Companies: &grpcadapter.CompanyClient{Client: pb.NewCompanyServiceClient(companyConn)}, Scorer: &grpcadapter.ScoringClient{Client: pb.NewScoringServiceClient(scoringConn)}, Storage: storage, External: &external.Client{BaseURL: platform.Env("EXTERNAL_MOCKS_URL", "http://localhost:8090"), HTTP: &http.Client{}}}
	run, cancel := context.WithCancel(ctx)
	defer cancel()
	go worker(run, secret, service)
	server := platform.NewGRPC(secret)
	pb.RegisterApplicationServiceServer(server, &grpcadapter.Server{Service: service})
	platform.Serve(server, platform.Env("GRPC_ADDR", ":9092"), cancel)
}

func worker(ctx context.Context, secret string, s *usecase.Service) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	n := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n++
			jobCtx, cancel := platform.SystemContext(secret)
			if e := s.Process(jobCtx); e != nil {
				slog.Error("worker failed", "error", e)
			}
			cancel()
			if n%60 == 0 {
				gc, cancel := context.WithTimeout(ctx, 20*time.Second)
				if e := s.Reconcile(gc); e != nil {
					slog.Error("object reconciliation failed", "error", e)
				}

				cancel()
			}
		}
	}
}
