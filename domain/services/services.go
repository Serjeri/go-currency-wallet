package services

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/Serjeri/proto-exchange/exchange"
	"gw-currency-wallet/domain/handlers"
	"gw-currency-wallet/domain/models"
	"gw-currency-wallet/domain/services/auth"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	Get(ctx context.Context, user *models.Login) (int, error)
	GetBalance(ctx context.Context, id int) (*models.Balance, error)
	UpdateBalance(ctx context.Context, id int, updateBalance *models.UpdateBalance, newAmount int) error
}

type UserService struct {
	repo UserRepository
	grpc pb.ExchangeServiceClient
}

func NewUserService(repo UserRepository, grpc pb.ExchangeServiceClient) *UserService {
	return &UserService{repo: repo, grpc: grpc}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	hashedPassword := handlers.HashedPassword(user.Password)
	user.Password = hashedPassword

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := auth.CreateToken(id)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *UserService) GetUser(ctx context.Context, user *models.Login) (string, error) {
	hashedPassword := handlers.HashedPassword(user.Password)
	user.Password = hashedPassword

	id, err := s.repo.Get(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	token, err := auth.CreateToken(id)
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

	var newAmount int
	switch updateBalance.Status {
	case "deposit":
		switch updateBalance.Currency {
		case "USD":
			newAmount = currentBalance.USD + updateBalance.Amount
			currentBalance.USD = newAmount
		case "RUB":
			newAmount = currentBalance.RUB + updateBalance.Amount
			currentBalance.RUB = newAmount
		case "EUR":
			newAmount = currentBalance.EUR + updateBalance.Amount
			currentBalance.EUR = newAmount
		}
	case "withdrawal":
		switch updateBalance.Currency {
		case "USD":
			if currentBalance.USD < updateBalance.Amount {
				return nil, errors.New("insufficient funds")
			}
			newAmount = currentBalance.USD - updateBalance.Amount
			currentBalance.USD = newAmount
		case "RUB":
			newAmount = currentBalance.RUB - updateBalance.Amount
			if newAmount < 0 {
				return nil, errors.New("insufficient funds")
			}
			currentBalance.RUB = newAmount
		case "EUR":
			newAmount = currentBalance.EUR - updateBalance.Amount
			if newAmount < 0 {
				return nil, errors.New("insufficient funds")
			}
			currentBalance.EUR = newAmount
		}
	default:
		return nil, fmt.Errorf("unknown status: %s", updateBalance.Status)
	}

	if err := s.repo.UpdateBalance(ctx, id, updateBalance, newAmount); err != nil {
		return nil, fmt.Errorf("failed to update %s balance: %w", updateBalance.Currency, err)
	}

	return currentBalance, nil
}

func (s *UserService) GetRates(ctx context.Context) (*pb.ExchangeRatesResponse, error) {
	return s.grpc.GetExchangeRates(ctx, &pb.Empty{})
}
