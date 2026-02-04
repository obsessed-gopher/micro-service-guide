package usecases

import (
	"context"
	"testing"

	"github.com/obsessed-gopher/micro-service-guide/internal/models"
	"github.com/obsessed-gopher/micro-service-guide/internal/types"
)

type mockRepository struct {
	users map[string]*models.User
}

func newMockRepository() *mockRepository {
	return &mockRepository{users: make(map[string]*models.User)}
}

func (m *mockRepository) Create(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockRepository) Find(ctx context.Context, filter models.UserFilter, pagination *models.Pagination) ([]*models.User, error) {
	var result []*models.User
	for _, user := range m.users {
		if !m.matchesFilter(user, filter) {
			continue
		}
		result = append(result, user)
	}
	if pagination != nil && pagination.Limit > 0 {
		start := min(pagination.Offset, len(result))
		end := min(start+pagination.Limit, len(result))
		result = result[start:end]
	}
	return result, nil
}

func (m *mockRepository) Count(ctx context.Context, filter models.UserFilter) (int, error) {
	count := 0
	for _, user := range m.users {
		if m.matchesFilter(user, filter) {
			count++
		}
	}
	return count, nil
}

func (m *mockRepository) matchesFilter(user *models.User, filter models.UserFilter) bool {
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

func (m *mockRepository) Update(ctx context.Context, user *models.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return types.ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, filter models.UserFilter) (int, error) {
	var toDelete []string
	for id, user := range m.users {
		if m.matchesFilter(user, filter) {
			toDelete = append(toDelete, id)
		}
	}
	for _, id := range toDelete {
		delete(m.users, id)
	}
	return len(toDelete), nil
}

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *mockHasher) Compare(hash, password string) bool {
	return hash == "hashed_"+password
}

type mockIDGen struct {
	counter int
}

func (m *mockIDGen) Generate() string {
	m.counter++
	return "test-id-" + string(rune('0'+m.counter))
}

func TestUserUsecase_Create(t *testing.T) {
	repo := newMockRepository()
	hasher := &mockHasher{}
	idGen := &mockIDGen{}
	usecase := NewUserUsecase(repo, hasher, idGen)

	tests := []struct {
		name    string
		input   models.CreateUserInput
		wantErr error
	}{
		{
			name: "valid user",
			input: models.CreateUserInput{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "invalid email",
			input: models.CreateUserInput{
				Email:    "invalid-email",
				Name:     "Test User",
				Password: "password123",
			},
			wantErr: types.ErrInvalidEmail,
		},
		{
			name: "short password",
			input: models.CreateUserInput{
				Email:    "test2@example.com",
				Name:     "Test User",
				Password: "short",
			},
			wantErr: types.ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := usecase.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}

			if user.Email != tt.input.Email {
				t.Errorf("Create() email = %v, want %v", user.Email, tt.input.Email)
			}

			if user.Status != types.UserStatusActive {
				t.Errorf("Create() status = %v, want %v", user.Status, types.UserStatusActive)
			}
		})
	}
}
