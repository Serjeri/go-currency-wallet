package query

import (
	"context"
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
		return 0, fmt.Errorf("username or password is incorrect")
	}

	return id, nil
}

func (r *UserRepository) GetBalance(ctx context.Context, id int) (*models.Balance, error) {
	var balance models.Balance
	err := r.client.QueryRow(ctx, `SELECT usd, rub, eur FROM wallet WHERE user_id = $1`, id).Scan(
		&balance.USD, &balance.RUB, &balance.EUR)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	return &balance, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id int, updateBalance *models.UpdateBalance, newAmount int) error {
	query := fmt.Sprintf("UPDATE wallet SET %s = $1 WHERE user_id = $2", updateBalance.Currency)

    result, err := r.client.Exec(ctx, query, newAmount, id)
    if err != nil {
        return fmt.Errorf("failed to update %s balance for user %d: %w",
            updateBalance.Currency, id, err)
    }

    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("no wallet found for user %d", id)
    }

    return nil
}
