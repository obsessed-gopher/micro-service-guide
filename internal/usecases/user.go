// Package usecases содержит бизнес-логику сервиса.
package usecases

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
	Find(ctx context.Context, filter models.UserFilter, pagination *models.Pagination) ([]*models.User, error)
	Count(ctx context.Context, filter models.UserFilter) (int, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, filter models.UserFilter) (int, error)
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

// UserUsecase - модуль бизнес-логики пользователей.
type UserUsecase struct {
	repo   UserRepository
	hasher PasswordHasher
	idGen  IDGenerator
}

// NewUserUsecase создаёт новый модуль пользователей.
func NewUserUsecase(repo UserRepository, hasher PasswordHasher, idGen IDGenerator) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		hasher: hasher,
		idGen:  idGen,
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// findOne возвращает одного пользователя по фильтру или ErrUserNotFound.
func (m *UserUsecase) findOne(ctx context.Context, filter models.UserFilter) (*models.User, error) {
	users, err := m.repo.Find(ctx, filter, &models.Pagination{Limit: 1})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, types.ErrUserNotFound
	}

	return users[0], nil
}

// Create создаёт нового пользователя.
func (m *UserUsecase) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	if !emailRegex.MatchString(input.Email) {
		return nil, types.ErrInvalidEmail
	}

	if len(input.Password) < 8 {
		return nil, types.ErrInvalidPassword
	}

	existing, err := m.findOne(ctx, models.UserFilter{Emails: []string{input.Email}})
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
func (m *UserUsecase) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := m.findOne(ctx, models.UserFilter{IDs: []string{id}})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

// Update обновляет данные пользователя.
func (m *UserUsecase) Update(ctx context.Context, id string, input models.UpdateUserInput) (*models.User, error) {
	user, err := m.findOne(ctx, models.UserFilter{IDs: []string{id}})
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

// Delete удаляет пользователей по фильтру. Возвращает количество удалённых.
func (m *UserUsecase) Delete(ctx context.Context, filter models.UserFilter) (int, error) {
	count, err := m.repo.Delete(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("delete users: %w", err)
	}
	return count, nil
}

// ListFilter - параметры для метода List (публичный API usecase).
type ListFilter struct {
	Filters models.UserFilter
	Limit   int
	Offset  int
}

// List возвращает список пользователей.
func (m *UserUsecase) List(ctx context.Context, filter ListFilter) ([]*models.User, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	repoFilter := models.UserFilter{
		IDs:      filter.Filters.IDs,
		Emails:   filter.Filters.Emails,
		Statuses: filter.Filters.Statuses,
	}

	pagination := &models.Pagination{Limit: filter.Limit, Offset: filter.Offset}

	users, err := m.repo.Find(ctx, repoFilter, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	total, err := m.repo.Count(ctx, repoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	return users, total, nil
}
