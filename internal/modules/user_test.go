package modules

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

func (m *mockRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, types.ErrUserNotFound
	}
	return user, nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, types.ErrUserNotFound
}

func (m *mockRepository) Update(ctx context.Context, user *models.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return types.ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.users[id]; !ok {
		return types.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *mockRepository) List(ctx context.Context, filter models.ListUsersFilter) ([]*models.User, int, error) {
	var users []*models.User
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, len(users), nil
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

func TestUserModule_Create(t *testing.T) {
	repo := newMockRepository()
	hasher := &mockHasher{}
	idGen := &mockIDGen{}
	module := NewUserModule(repo, hasher, idGen)

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
			user, err := module.Create(context.Background(), tt.input)

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
