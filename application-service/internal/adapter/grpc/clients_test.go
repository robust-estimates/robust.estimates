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
