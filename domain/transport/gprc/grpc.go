package gprc

import (
	"fmt"
	"net"

	pb "github.com/Serjeri/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// type Client struct {
// 	grpc pb.ExchangeServiceClient
// 	conn *grpc.ClientConn
// }


func New(host, port string) (pb.ExchangeServiceClient, func() error) {
	// conn, err := grpc.NewClient(net.JoinHostPort(host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	fmt.Errorf("Failed to create gRPC client %w", err)
	// }

	// serviceClient := pb.NewExchangeServiceClient(conn)
	// client := &Client{
	// 	grpc: serviceClient,
	// 	conn: conn,
	// }
	// return client, conn.Close

	conn, err := grpc.NewClient(net.JoinHostPort(host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("Failed to create gRPC client %w", err)
	}

	client := pb.NewExchangeServiceClient(conn)
	return client, conn.Close
}
