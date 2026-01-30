// Package metrics содержит бизнес-метрики сервиса.
package metrics

import (
	"sync/atomic"
)

// UserMetrics - метрики, связанные с пользователями.
type UserMetrics struct {
	usersCreated atomic.Int64
	usersDeleted atomic.Int64
	usersBlocked atomic.Int64
}

// NewUserMetrics создаёт новые метрики пользователей.
func NewUserMetrics() *UserMetrics {
	return &UserMetrics{}
}

// IncUsersCreated увеличивает счётчик созданных пользователей.
func (m *UserMetrics) IncUsersCreated() {
	m.usersCreated.Add(1)
}

// IncUsersDeleted увеличивает счётчик удалённых пользователей.
func (m *UserMetrics) IncUsersDeleted() {
	m.usersDeleted.Add(1)
}

// IncUsersBlocked увеличивает счётчик заблокированных пользователей.
func (m *UserMetrics) IncUsersBlocked() {
	m.usersBlocked.Add(1)
}

// Stats возвращает текущие значения метрик.
func (m *UserMetrics) Stats() (created, deleted, blocked int64) {
	return m.usersCreated.Load(), m.usersDeleted.Load(), m.usersBlocked.Load()
}