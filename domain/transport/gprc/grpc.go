package gprc

import (
	"time"

	"github.com/gofiber/fiber/v2/log"

	pb "github.com/Serjeri/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(address string) (pb.ExchangeServiceClient, func() error) {
    conn, err := grpc.NewClient(
        address, // Используем переданный адрес
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(), // Опционально: ждём подключения
        grpc.WithTimeout(5*time.Second), // Таймаут подключения
    )
    if err != nil {
        log.Fatalf("failed to create gRPC client to %s: %v", address, err)
    }

    client := pb.NewExchangeServiceClient(conn)
    return client, conn.Close
}
