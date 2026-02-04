package repository

import (
	"context"
	"fmt"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	"github.com/obsessed-gopher/micro-service-guide/internal/types"
)

// Create сохраняет пользователя в БД.
func (r *PostgresRepository) Create(ctx context.Context, user *models.User) error {
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

// Find возвращает пользователей по фильтру.
func (r *PostgresRepository) Find(ctx context.Context, filter models.UserFilter, pagination *models.Pagination) ([]*models.User, error) {
	qb := newQueryBuilder()
	qb.buildUserFilter(filter)

	query := `SELECT id, email, name, password_hash, status, created_at, updated_at FROM users` +
		qb.whereClause() +
		` ORDER BY created_at DESC` +
		qb.addPagination(pagination)

	rows, err := r.db.QueryContext(ctx, query, qb.args...)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.PasswordHash,
			&user.Status, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

// Count возвращает количество пользователей по фильтру.
func (r *PostgresRepository) Count(ctx context.Context, filter models.UserFilter) (int, error) {
	qb := newQueryBuilder()
	qb.buildUserFilter(filter)

	query := `SELECT COUNT(*) FROM users` + qb.whereClause()

	var count int

	if err := r.db.QueryRowContext(ctx, query, qb.args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}

// Update обновляет пользователя в БД.
func (r *PostgresRepository) Update(ctx context.Context, user *models.User) error {
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

// Delete удаляет пользователей по фильтру. Возвращает количество удалённых.
func (r *PostgresRepository) Delete(ctx context.Context, filter models.UserFilter) (int, error) {
	qb := newQueryBuilder()
	qb.buildUserFilter(filter)

	query := `DELETE FROM users` + qb.whereClause()

	result, err := r.db.ExecContext(ctx, query, qb.args...)
	if err != nil {
		return 0, fmt.Errorf("delete users: %w", err)
	}

	count, _ := result.RowsAffected()

	return int(count), nil
}

// buildUserFilter применяет фильтр пользователей к query builder.
func (qb *queryBuilder) buildUserFilter(filter models.UserFilter) {
	if len(filter.IDs) > 0 {
		qb.addInCondition("id", toAnySlice(filter.IDs))
	}

	if len(filter.Emails) > 0 {
		qb.addInCondition("email", toAnySlice(filter.Emails))
	}

	if len(filter.Statuses) > 0 {
		qb.addInCondition("status", statusesToAny(filter.Statuses))
	}
}

func statusesToAny(statuses []types.UserStatus) []any {
	result := make([]any, len(statuses))

	for i, s := range statuses {
		result[i] = s
	}

	return result
}
