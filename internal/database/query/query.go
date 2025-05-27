package query

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/database"
	"gw-currency-wallet/internal/models"
)

type Repository struct {
	client database.Client
}

func NewRepository(client database.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) RegistrUser(ctx context.Context, user models.User) (int, error) {
	var id int
	err := r.client.QueryRow(
		ctx,
		`INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id`,
		user.Name, user.Password, user.Email,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	return id, nil
}

// func (r *Repository) GetAllTasks(ctx context.Context) (u []models.Tasks, err error) {
// 	rows, err := r.client.Query(ctx, `SELECT * FROM tasks`)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tasks := make([]models.Tasks, 0)

// 	for rows.Next() {
// 		var task models.Tasks

// 		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Created_at, &task.Updated_at)
// 		if err != nil {
// 			return nil, err
// 		}

// 		tasks = append(tasks, task)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return tasks, nil
// }

// func (r *Repository) UpdateTask(ctx context.Context, id int, title string, description string) (bool, error) {
// 	var exists bool
// 	err := r.client.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`, id).Scan(&exists)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check task existence: %w", err)
// 	}

// 	if !exists {
// 		return false, nil
// 	}
// 	_, err = r.client.Exec(ctx, `UPDATE tasks SET title = $1, description = $2, status='in_progress' WHERE id = $3`, title, description, id)
// 	if err != nil {
// 		return false, fmt.Errorf("Ошибка обновления таски: %w", err)
// 	}

// 	return true, nil
// }

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
