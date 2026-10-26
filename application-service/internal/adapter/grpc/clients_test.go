package grpcadapter

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"robust/application-service/internal/domain"
	pb "robust/contracts/credit/v1"
	"robust/platform"
)

type recoveringScoring struct {
	pb.UnimplementedScoringServiceServer
}

func (recoveringScoring) Evaluate(_ context.Context, in *pb.ScoringInput) (*pb.ScoringOutput, error) {
	return &pb.ScoringOutput{
		ScoringRequestId:   in.ScoringRequestId,
		ApplicationId:      in.ApplicationId,
		ApplicationVersion: in.ApplicationVersion,
		Scoring:            &pb.Scoring{Score: 740, RiskLevel: "LOW", ModelVersion: "rules-v1", Recommendation: "APPROVE"},
	}, nil
}

func TestScoringWaitsForReconnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	address := listener.Addr().String()
	listener.Close()
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	connection.Connect()
	for state := connection.GetState(); state != connectivity.TransientFailure; state = connection.GetState() {
		if !connection.WaitForStateChange(ctx, state) {
			t.Fatal("client did not observe the unavailable scoring service")
		}
	}

	listener, err = net.Listen("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterScoringServiceServer(server, recoveringScoring{})
	defer server.Stop()
	go server.Serve(listener)

	client := ScoringClient{Client: pb.NewScoringServiceClient(connection)}
	result, err := client.Evaluate(ctx, domain.ScoringInput{ScoringRequestID: "request", ApplicationID: "application", ApplicationVersion: 1})
	if err != nil {
		t.Fatalf("scoring should reconnect after temporary unavailability: %v", err)
	}
	if result.Score != 740 || result.RiskLevel != "LOW" {
		t.Fatalf("unexpected recovered result: %+v", result)
	}
}

func TestFeaturesKeepExternalDataInProto(t *testing.T) {
	in := domain.ScoringInput{ScoringRequestID: "job", ApplicationID: "app", ApplicationVersion: 1, Features: domain.Features{
		Registry:     &domain.RegistryData{Status: "ACTIVE", OkvedSection: "G", TaxDebtMinor: 500},
		CreditReport: &domain.CreditReport{Owner: &domain.CreditHistory{Found: true, OverdueCount12m: 2, MaxOverdueDays24m: 45}},
	}}
	out := platform.PB(in, &pb.ScoringInput{})
	f := out.Features
	if f.Registry.GetOkvedSection() != "G" || f.Registry.GetTaxDebtMinor() != 500 || f.CreditReport.GetOwner().GetOverdueCount_12M() != 2 {
		t.Fatalf("external data lost in proto: %v", f)
	}
	if f.CreditReport.Company != nil || f.Account != nil {
		t.Fatal("missing sources must stay unset")
	}
}
