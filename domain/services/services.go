package services

import (
	"context"
	"fmt"

	"gw-currency-wallet/domain/broker/kafka"
	"gw-currency-wallet/domain/lib"
	"gw-currency-wallet/domain/lib/jwttoken"
	"gw-currency-wallet/domain/models"

	pb "github.com/Serjeri/proto-exchange/exchange"
	"github.com/gofiber/fiber/v2/log"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	Get(ctx context.Context, user *models.Login) (int, error)
	GetBalance(ctx context.Context, id int) (*models.Balance, error)
	UpdateBalance(ctx context.Context, id int, updateBalance *models.UpdateBalance, newAmount int) error
	UpdateBalanceExchange(ctx context.Context, id int, FromCurrency, ToCurrency string, FromAmount, ToAmount int) error
}

type UserService struct {
	repo UserRepository
	grpc pb.ExchangeServiceClient
}

func NewUserService(repo UserRepository, grpc pb.ExchangeServiceClient) *UserService {
	return &UserService{repo: repo, grpc: grpc}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	hashedPassword := lib.HashedPassword(user.Password)
	user.Password = hashedPassword

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := jwttoken.CreateToken(id)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *UserService) GetUser(ctx context.Context, user *models.Login) (string, error) {
	hashedPassword := lib.HashedPassword(user.Password)
	user.Password = hashedPassword

	id, err := s.repo.Get(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	token, err := jwttoken.CreateToken(id)
	if err != nil {
		return "", fmt.Errorf("failed to generate auth token: %w", err)
	}

	return token, nil
}

func (s *UserService) GetBalanceUser(ctx context.Context, id int) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func (s *UserService) UpdateBalanceUser(ctx context.Context, id int, updateBalance *models.UpdateBalance) (*models.Balance, error) {
	currentBalance, err := s.repo.GetBalance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current balance: %w", err)
	}

	newAmount, err := lib.UpdateBalance(currentBalance, updateBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	if err := s.repo.UpdateBalance(ctx, id, updateBalance, newAmount); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	return currentBalance, nil
}

func (s *UserService) GetRates(ctx context.Context) (*pb.ExchangeRatesResponse, error) {
	return s.grpc.GetExchangeRates(ctx, &pb.Empty{})
}

func (s *UserService) Exchange(ctx context.Context, user *models.Exchange, id int) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	newAmount, err := lib.DeductFromBalance(balance, user.FromCurrency, user.Amount)
	if err != nil {
		log.Info("failed to change balace: %w", err)
		return nil, fmt.Errorf("failed to change balace: %w", err)
	}

	exchange, err := s.grpc.PerformExchange(ctx, &pb.ExchangeRequest{
		Amount: int64(user.Amount), FromCurrency: user.FromCurrency, ToCurrency: user.ToCurrency,
	})
	if err != nil {
		log.Info("failed to perform exchange: %w", err)
		return nil, fmt.Errorf("failed to perform exchange: %w", err)
	}

	exchangedAmount := int(exchange.NewBalance[user.ToCurrency] * 10000)

	toAmount := lib.AddToBalance(exchangedAmount, user.ToCurrency, balance)

	if user.Amount >= 30000 {
		success, err := kafka.Producer(id, user)
		if err != nil {
			return nil, fmt.Errorf("failed to send exchange request to kafka: %w", err)
		}

		if !success {
			return nil, fmt.Errorf("kafka producer failed to send message")
		}
	}

	if err := s.repo.UpdateBalanceExchange(ctx, id, user.FromCurrency, user.ToCurrency, newAmount, toAmount); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	return balance, nil
}
