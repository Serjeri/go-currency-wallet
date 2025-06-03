package services

import (
	"context"
	"errors"
	"fmt"
	"gw-currency-wallet/domain/handlers"
	"gw-currency-wallet/domain/models"
	"gw-currency-wallet/domain/services/auth"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	Get(ctx context.Context, user *models.Login) (int, error)
	CheckUser(ctx context.Context, id int) (bool, error)
	GetBalance(ctx context.Context, id int) (*models.Balance, error)
	UpdateBalance(ctx context.Context, id int, updateBalance *models.UpdateBalance, newAmount int) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	hashedPassword := handlers.HashedPassword(user.Password)
	user.Password = hashedPassword

	id, err := s.repo.Create(context.TODO(), user)
	if err != nil {
		err := errors.New("user with this name or email already exists")
		return err.Error(), nil
	}

	token, err := auth.CreateToken(id)
	if err != nil {
		err := errors.New("token generation failed")
		return err.Error(), nil
	}

	return token, nil
}

func (s *UserService) GetUser(ctx context.Context, user *models.Login) (string, error) {
	hashedPassword := handlers.HashedPassword(user.Password)
	user.Password = hashedPassword

	id, err := s.repo.Get(context.TODO(), user)
	if err != nil {
		err := errors.New("unauthorized")
		return err.Error(), nil
	}

	token, err := auth.CreateToken(id)
	if err != nil {
		err := errors.New("token generation failed")
		return err.Error(), nil
	}

	return token, nil
}

func (s *UserService) GetBalanceUser(ctx context.Context, id int) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(context.TODO(), id)
	if err != nil {
		err := errors.New("token generation failed")
		return nil, err
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

    err = s.repo.UpdateBalance(ctx, id, updateBalance, newAmount)
    if err != nil {
        return nil, fmt.Errorf("failed to update balance: %w", err)
    }

	return currentBalance, nil
}
