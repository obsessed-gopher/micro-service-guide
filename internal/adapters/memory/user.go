// Package memory содержит реализации репозиториев.
package memory

import (
	"context"
	"sync"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	"github.com/obsessed-gopher/micro-service-guide/internal/types"
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

// Find возвращает пользователей по фильтру.
func (r *MemoryUserRepository) Find(ctx context.Context, filter models.UserFilter, pagination *models.Pagination) ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*models.User

	for _, user := range r.users {
		if !r.matchesFilter(user, filter) {
			continue
		}

		filtered = append(filtered, user)
	}

	if pagination == nil {
		return filtered, nil
	}

	start := min(pagination.Offset, len(filtered))
	end := min(start+pagination.Limit, len(filtered))

	return filtered[start:end], nil
}

// Count возвращает количество пользователей по фильтру.
func (r *MemoryUserRepository) Count(ctx context.Context, filter models.UserFilter) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0

	for _, user := range r.users {
		if r.matchesFilter(user, filter) {
			count++
		}
	}

	return count, nil
}

// matchesFilter проверяет, соответствует ли пользователь фильтру.
func (r *MemoryUserRepository) matchesFilter(user *models.User, filter models.UserFilter) bool {
	if len(filter.IDs) > 0 && !containsString(filter.IDs, user.ID) {
		return false
	}

	if len(filter.Emails) > 0 && !containsString(filter.Emails, user.Email) {
		return false
	}

	if len(filter.Statuses) > 0 && !containsStatus(filter.Statuses, user.Status) {
		return false
	}

	return true
}

func containsString(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}

	return false
}

func containsStatus(slice []types.UserStatus, val types.UserStatus) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}

	return false
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

// Delete удаляет пользователей по фильтру. Возвращает количество удалённых.
func (r *MemoryUserRepository) Delete(ctx context.Context, filter models.UserFilter) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var toDelete []string

	for id, user := range r.users {
		if r.matchesFilter(user, filter) {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		delete(r.users, id)
	}

	return len(toDelete), nil
}
