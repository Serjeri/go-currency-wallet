package query

import (
	"context"
	"errors"
	"fmt"
	"gw-currency-wallet/domain/models"
	"gw-currency-wallet/domain/repository"

	"strings"
)

type UserRepository struct {
	client repository.Client
}

func NewRepository(client repository.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var exists bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE name = $1 OR email = $2)`,
		strings.ToUpper(user.Name),
		strings.ToUpper(user.Email),
	).Scan(&exists)

	if err != nil {
		return 0, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return 0, fmt.Errorf("user with this name or email already exists")
	}

	var id int
	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`,
		strings.ToUpper(user.Name),
		user.Password,
		strings.ToUpper(user.Email),
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO wallet (user_id, usd, rub, eur) VALUES ($1, 0.0, 0.0, 0.0)`,
		id,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create wallet: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (r *UserRepository) Get(ctx context.Context, user *models.Login) (int, error) {
	var dbPassword string
	var id int

	r.client.QueryRow(ctx, `SELECT id, password FROM users WHERE name = $1`,
		strings.ToUpper(user.Name)).Scan(&id, &dbPassword)

	if dbPassword != user.Password {
		return 0, errors.New("username or password is incorrect")
	}

	return id, nil
}

func (r *UserRepository) CheckUser(ctx context.Context, id int) (bool, error) {
	var exists bool

	err := r.client.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to create user: %w", err)
	}

	if !exists {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) GetBalance(ctx context.Context, id int) (*models.Balance, error) {
	var balance models.Balance
	err := r.client.QueryRow(ctx, `SELECT usd, rub, eur FROM wallet WHERE user_id = $1`, id).Scan(
		&balance.USD, &balance.RUB, &balance.EUR)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &balance, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id int, updateBalance *models.UpdateBalance) (*models.Balance, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentBalance models.Balance
	err = tx.QueryRow(ctx,
		`SELECT usd, rub, eur FROM wallet WHERE user_id = $1`, id).
		Scan(&currentBalance.USD, &currentBalance.RUB, &currentBalance.EUR)
	if err != nil {
		return nil, fmt.Errorf("failed to get current balance: %w", err)
	}

	var newAmount float64
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

	query := fmt.Sprintf("UPDATE wallet SET %s = $1 WHERE user_id = $2", updateBalance.Currency)
    _, err = tx.Exec(ctx, query, newAmount, id)
    if err != nil {
        return nil, fmt.Errorf("failed to update balance: %w", err)
    }

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &currentBalance, nil
}

// func (r *Repository) DeleteTask(ctx context.Context, id int) (bool, error) {
// 	var exists bool
// 	err := r.client.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`, id).Scan(&exists)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check task existence: %w", err)
// 	}

// 	if !exists {
// 		return false, nil
// 	}

// 	_, err = r.client.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
