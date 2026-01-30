// Package modules содержит бизнес-логику сервиса.
package modules

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	"github.com/obsessed-gopher/micro-service-guide/internal/types"
)

// UserRepository - интерфейс репозитория пользователей.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter models.ListUsersFilter) ([]*models.User, int, error)
}

// PasswordHasher - интерфейс для хэширования паролей.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// IDGenerator - интерфейс генератора ID.
type IDGenerator interface {
	Generate() string
}

// UserModule - модуль бизнес-логики пользователей.
type UserModule struct {
	repo   UserRepository
	hasher PasswordHasher
	idGen  IDGenerator
}

// NewUserModule создаёт новый модуль пользователей.
func NewUserModule(repo UserRepository, hasher PasswordHasher, idGen IDGenerator) *UserModule {
	return &UserModule{
		repo:   repo,
		hasher: hasher,
		idGen:  idGen,
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Create создаёт нового пользователя.
func (m *UserModule) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	if !emailRegex.MatchString(input.Email) {
		return nil, types.ErrInvalidEmail
	}

	if len(input.Password) < 8 {
		return nil, types.ErrInvalidPassword
	}

	existing, err := m.repo.GetByEmail(ctx, input.Email)
	if err != nil && !types.IsNotFound(err) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, types.ErrUserAlreadyExists
	}

	hash, err := m.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := &models.User{
		ID:           m.idGen.Generate(),
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: hash,
		Status:       types.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := m.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

// GetByID возвращает пользователя по ID.
func (m *UserModule) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

// Update обновляет данные пользователя.
func (m *UserModule) Update(ctx context.Context, id string, input models.UpdateUserInput) (*models.User, error) {
	user, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if user.IsBlocked() {
		return nil, types.ErrUserBlocked
	}

	if input.Email != nil {
		if !emailRegex.MatchString(*input.Email) {
			return nil, types.ErrInvalidEmail
		}
		user.Email = *input.Email
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Status != nil {
		user.Status = *input.Status
	}

	user.UpdatedAt = time.Now()

	if err := m.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

// Delete удаляет пользователя.
func (m *UserModule) Delete(ctx context.Context, id string) error {
	if err := m.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// List возвращает список пользователей.
func (m *UserModule) List(ctx context.Context, filter models.ListUsersFilter) ([]*models.User, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	users, total, err := m.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return users, total, nil
}
