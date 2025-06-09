package gprc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/Serjeri/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(host, port string) (pb.ExchangeServiceClient, func() error) {
	conn, err := grpc.NewClient(net.JoinHostPort(host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", fmt.Errorf("connection error: %w", err))
	}

	client := pb.NewExchangeServiceClient(conn)
	return client, conn.Close
}
