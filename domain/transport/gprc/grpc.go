package gprc

import (
	"fmt"
	"log"

	pb "github.com/Serjeri/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(address string) (pb.ExchangeServiceClient, func() error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", fmt.Errorf("connection error: %w", err))
	}

	client := pb.NewExchangeServiceClient(conn)
	return client, conn.Close
}
