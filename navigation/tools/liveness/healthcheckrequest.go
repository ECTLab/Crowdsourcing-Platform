package health

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func PerformHealthCheckRequest(address string) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	healthRequest := &grpc_health_v1.HealthCheckRequest{
		Service: "",
	}

	_, err = client.Check(ctx, healthRequest)

	if err != nil {
		log.Fatalf("Health check request failed: %v", err)
	}
}
