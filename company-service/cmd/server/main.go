package main

import (
	"context"
	"os"

	"robust/company-service/internal/adapter/grpc"
	"robust/company-service/internal/adapter/postgres"
	"robust/company-service/internal/usecase"
	pb "robust/contracts/credit/v1"
	"robust/platform"
)

func main() {
	if platform.Healthcheck(platform.Env("GRPC_ADDR", ":9091")) {
		return
	}

	ctx := context.Background()
	db := platform.OpenDB(ctx, os.Getenv("DATABASE_URL"))
	defer db.Close()
	e := platform.Migrate(ctx, db, platform.Env("MIGRATIONS_DIR", "migrations"))
	if e != nil {
		panic(e)
	}

	r := &postgres.Repository{DB: db}
	if e = r.Seed(ctx); e != nil {
		panic(e)
	}

	s := platform.NewGRPC(platform.Secret())
	pb.RegisterCompanyServiceServer(s, &grpcadapter.Server{Service: &usecase.Service{Repo: r}})
	platform.Serve(s, platform.Env("GRPC_ADDR", ":9091"), func() {})
}
