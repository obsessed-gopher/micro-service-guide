package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/example/user-service/internal/models"
	"github.com/example/user-service/internal/types"
)

// PostgresUserRepository - PostgreSQL реализация репозитория.
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository создаёт новый PostgreSQL репозиторий.
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create сохраняет пользователя в БД.
func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Name, user.PasswordHash,
		user.Status, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

// GetByID возвращает пользователя по ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, status, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, types.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	return user, nil
}

// GetByEmail возвращает пользователя по email.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, status, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, types.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query user by email: %w", err)
	}
	return user, nil
}

// Update обновляет пользователя в БД.
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET email = $2, name = $3, status = $4, updated_at = $5
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Name, user.Status, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return types.ErrUserNotFound
	}
	return nil
}

// Delete удаляет пользователя из БД.
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return types.ErrUserNotFound
	}
	return nil
}

// List возвращает список пользователей.
func (r *PostgresUserRepository) List(ctx context.Context, filter models.ListUsersFilter) ([]*models.User, int, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	args := []any{}
	argNum := 1

	if filter.Status != nil {
		countQuery += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, *filter.Status)
		argNum++
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Select query
	query := `SELECT id, email, name, password_hash, status, created_at, updated_at FROM users WHERE 1=1`
	if filter.Status != nil {
		query += " AND status = $1"
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.PasswordHash,
			&user.Status, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, total, rows.Err()
}