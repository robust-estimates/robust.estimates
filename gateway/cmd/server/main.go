package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "robust/contracts/credit/v1"
	httpadapter "robust/gateway/internal/adapter/http"
	"robust/gateway/internal/adapter/postgres"
	"robust/gateway/internal/domain"
	"robust/gateway/internal/usecase"
	"robust/platform"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		c := http.Client{Timeout: 3 * time.Second}
		r, e := c.Get("http://127.0.0.1" + platform.Env("HTTP_ADDR", ":8080") + "/readyz")
		if e != nil {
			os.Exit(1)
		}

		r.Body.Close()
		if r.StatusCode != 200 {
			os.Exit(1)
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

	users := &postgres.Users{DB: db}
	if os.Getenv("SEED_DEMO") == "true" {
		if e = users.Seed(ctx, os.Getenv("DEMO_PASSWORD")); e != nil {
			panic(e)
		}
	}

	company := platform.Dial(platform.Env("COMPANY_GRPC_ADDR", "localhost:9091"))
	defer company.Close()
	applications := platform.Dial(platform.Env("APPLICATION_GRPC_ADDR", "localhost:9092"))
	defer applications.Close()
	secret := platform.Secret()
	api := &httpadapter.Server{Auth: &usecase.Auth{Users: users, Sign: func(u domain.User) (string, error) { return platform.Token(secret, u.ID, u.DisplayName, u.Role) }}, DB: db, CompanyConn: company, ApplicationConn: applications, Companies: pb.NewCompanyServiceClient(company), Applications: pb.NewApplicationServiceClient(applications), Secret: secret, Origin: platform.Env("CORS_ORIGIN", "http://localhost:5173"), OpenAPI: platform.Env("OPENAPI_FILE", "../api/public/openapi.yaml"), Swagger: platform.Env("SWAGGER_FILE", "../docs/swagger.html")}
	server := &http.Server{Addr: platform.Env("HTTP_ADDR", ":8080"), Handler: api.Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 70 * time.Second, WriteTimeout: 70 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16384}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()
	slog.Info("gateway listening", "address", server.Addr)
	if e = server.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
		panic(e)
	}
}
