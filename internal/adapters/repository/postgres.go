package repository

import (
	"database/sql"
)

// PostgresRepository - PostgreSQL реализация репозитория.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository создаёт новый PostgreSQL репозиторий.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
