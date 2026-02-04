// Package models содержит бизнес-модели сервиса.
package models

import (
	"time"

	"github.com/obsessed-gopher/micro-service-guide/internal/types"
)

// User - бизнес-модель пользователя.
type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	Status       types.UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// IsActive проверяет, активен ли пользователь.
func (u *User) IsActive() bool {
	return u.Status == types.UserStatusActive
}

// IsBlocked проверяет, заблокирован ли пользователь.
func (u *User) IsBlocked() bool {
	return u.Status == types.UserStatusBlocked
}

// CreateUserInput - входные данные для создания пользователя.
type CreateUserInput struct {
	Email    string
	Name     string
	Password string
}

// UpdateUserInput - входные данные для обновления пользователя.
type UpdateUserInput struct {
	Email  *string
	Name   *string
	Status *types.UserStatus
}

// UserFilter - фильтры для поиска пользователей.
// Все поля — слайсы для поддержки множественных значений (IN).
// Пустой слайс означает "без фильтра по этому полю".
type UserFilter struct {
	IDs      []string
	Emails   []string
	Statuses []types.UserStatus
}

// Pagination - параметры пагинации.
type Pagination struct {
	Limit  int
	Offset int
}
