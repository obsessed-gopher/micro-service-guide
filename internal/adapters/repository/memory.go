// Package repository содержит реализации репозиториев.
package repository

import (
	"context"
	"sync"

	"github.com/example/user-service/internal/models"
	"github.com/example/user-service/internal/types"
)

// MemoryUserRepository - in-memory реализация репозитория (для тестов и демо).
type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

// NewMemoryUserRepository создаёт новый in-memory репозиторий.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

// Create сохраняет пользователя.
func (r *MemoryUserRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

// GetByID возвращает пользователя по ID.
func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, types.ErrUserNotFound
	}
	return user, nil
}

// GetByEmail возвращает пользователя по email.
func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, types.ErrUserNotFound
}

// Update обновляет пользователя.
func (r *MemoryUserRepository) Update(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return types.ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}

// Delete удаляет пользователя.
func (r *MemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return types.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

// List возвращает список пользователей с фильтрацией.
func (r *MemoryUserRepository) List(ctx context.Context, filter models.ListUsersFilter) ([]*models.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*models.User
	for _, user := range r.users {
		if filter.Status != nil && user.Status != *filter.Status {
			continue
		}
		filtered = append(filtered, user)
	}

	total := len(filtered)

	start := min(filter.Offset, len(filtered))
	end := min(start+filter.Limit, len(filtered))

	return filtered[start:end], total, nil
}