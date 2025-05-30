package services

import (
	"context"
	"errors"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/services/auth"
	"gw-currency-wallet/internal/services/handlers"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	Get(ctx context.Context, user *models.Login) (int, error)
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
