package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (title, description, status, created_at, updated_at, recurrence_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, status, created_at, updated_at, recurrence_id
	`
	row := r.pool.QueryRow(ctx, query,
		task.Title, task.Description, task.Status,
		task.CreatedAt, task.UpdatedAt, task.RecurrenceID)
	return scanTask(row)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at, recurrence_id
		FROM tasks
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	task, err := scanTask(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, taskdomain.ErrNotFound
	}
	return task, err
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4, recurrence_id = $5
		WHERE id = $6
		RETURNING id, title, description, status, created_at, updated_at, recurrence_id
	`
	row := r.pool.QueryRow(ctx, query,
		task.Title, task.Description, task.Status,
		task.UpdatedAt, task.RecurrenceID, task.ID)
	updated, err := scanTask(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, taskdomain.ErrNotFound
	}
	return updated, err
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}
	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at, recurrence_id
		FROM tasks
		ORDER BY id DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []taskdomain.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

func scanTask(scanner interface {
	Scan(dest ...any) error
}) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
	)
	err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RecurrenceID,
	)
	if err != nil {
		return nil, err
	}
	task.Status = taskdomain.Status(status)
	return &task, nil
}