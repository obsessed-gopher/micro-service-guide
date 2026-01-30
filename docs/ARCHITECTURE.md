# Архитектура микросервиса

## Обзор слоёв

```
┌─────────────────────────────────────────────────────────────┐
│                      cmd/user_service                       │
│                    (точка входа, DI)                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 internal/app/grpc/user_service              │
│           (gRPC handlers, валидация, маппинг proto)         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     internal/modules                        │
│                   (бизнес-логика)                           │
└─────────────────────────────────────────────────────────────┘
                    │                   │
                    ▼                   ▼
┌────────────────────────┐   ┌────────────────────────────────┐
│   internal/models      │   │      internal/adapters         │
│  (бизнес-сущности)     │   │   (БД, внешние API, брокеры)   │
└────────────────────────┘   └────────────────────────────────┘
                    │                   │
                    └─────────┬─────────┘
                              ▼
                    ┌─────────────────┐
                    │ internal/types  │
                    │ (ошибки, enum)  │
                    └─────────────────┘
```

## Структура Proto файлов

```
api/user_service/
├── api.proto               # service UserService { rpc... }
├── enum.proto              # enum UserStatus { ... }
├── model.proto             # message User { ... }
├── rpc_create_user.proto   # CreateUserRequest, CreateUserResponse
├── rpc_get_user.proto      # GetUserRequest, GetUserResponse
├── rpc_update_user.proto   # UpdateUserRequest, UpdateUserResponse
├── rpc_delete_user.proto   # DeleteUserRequest, DeleteUserResponse
└── rpc_list_users.proto    # ListUsersRequest, ListUsersResponse
```

**Правила:**
- Каждый RPC метод — отдельный файл `rpc_<method>.proto`
- Enum'ы в `enum.proto`
- Общие модели в `model.proto`
- Сервис в `api.proto` (импортирует остальные)

## Структура Handlers

```
internal/app/grpc/user_service/
├── server.go           # Server struct, интерфейсы, NewServer()
├── create_user.go      # func (s *Server) CreateUser(...)
├── get_user.go         # func (s *Server) GetUser(...)
├── update_user.go      # func (s *Server) UpdateUser(...)
├── delete_user.go      # func (s *Server) DeleteUser(...)
├── list_users.go       # func (s *Server) ListUsers(...)
├── converter.go        # userToProto(), statusFromProto()...
└── errors.go           # mapError()
```

**Правила:**
- Каждый RPC метод — отдельный файл `<method>.go`
- Конвертеры proto ↔ models в `converter.go`
- Маппинг ошибок в `errors.go`

## Детальное описание слоёв

### cmd/ — Точка входа

```go
func main() {
    cfg := config.Load(...)

    // Adapters
    repo := repository.NewPostgresUserRepository(db)
    hasher := hasher.NewBcryptHasher(0)

    // Business logic
    userModule := modules.NewUserModule(repo, hasher, idGen)

    // Transport
    server := userservice.NewServer(userModule)

    // gRPC
    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, server)
    grpcServer.Serve(listener)
}
```

### internal/app/grpc/ — Transport Layer

**server.go:**
```go
type UserModule interface {
    Create(ctx context.Context, input models.CreateUserInput) (*models.User, error)
    GetByID(ctx context.Context, id string) (*models.User, error)
    // ...
}

type Server struct {
    pb.UnimplementedUserServiceServer
    userModule UserModule
}

func NewServer(userModule UserModule) *Server {
    return &Server{userModule: userModule}
}
```

**create_user.go:**
```go
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    // 1. Валидация
    if req.Email == "" {
        return nil, status.Error(codes.InvalidArgument, "email required")
    }

    // 2. Вызов бизнес-логики
    user, err := s.userModule.Create(ctx, models.CreateUserInput{
        Email: req.Email,
        Name:  req.Name,
    })
    if err != nil {
        return nil, mapError(err)
    }

    // 3. Маппинг в proto
    return &pb.CreateUserResponse{User: userToProto(user)}, nil
}
```

**converter.go:**
```go
func userToProto(u *models.User) *pb.User {
    return &pb.User{
        Id:    u.ID,
        Email: u.Email,
        // ...
    }
}

func statusFromProto(s pb.UserStatus) types.UserStatus {
    switch s {
    case pb.UserStatus_USER_STATUS_ACTIVE:
        return types.UserStatusActive
    // ...
    }
}
```

**errors.go:**
```go
func mapError(err error) error {
    switch err {
    case types.ErrUserNotFound:
        return status.Error(codes.NotFound, err.Error())
    case types.ErrUserAlreadyExists:
        return status.Error(codes.AlreadyExists, err.Error())
    default:
        return status.Error(codes.Internal, "internal error")
    }
}
```

### internal/modules/ — Business Logic

```go
// Интерфейсы определяются здесь (consumer-side)
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id string) (*models.User, error)
}

type UserModule struct {
    repo   UserRepository
    hasher PasswordHasher
    idGen  IDGenerator
}

func (m *UserModule) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
    // Бизнес-валидация
    if !isValidEmail(input.Email) {
        return nil, types.ErrInvalidEmail
    }

    // Бизнес-логика
    user := &models.User{
        ID:     m.idGen.Generate(),
        Email:  input.Email,
        Status: types.UserStatusActive,
    }

    if err := m.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}
```

### internal/models/ — Domain Models

```go
type User struct {
    ID           string
    Email        string
    Name         string
    PasswordHash string
    Status       types.UserStatus
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (u *User) IsActive() bool {
    return u.Status == types.UserStatusActive
}
```

### internal/adapters/ — Infrastructure

```go
// repository/postgres.go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO users (id, email, ...) VALUES ($1, $2, ...)`,
        user.ID, user.Email, ...
    )
    return err
}
```

### internal/types/ — Types & Errors

```go
// errors.go
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidEmail      = errors.New("invalid email")
)

// status.go
type UserStatus int

const (
    UserStatusUnspecified UserStatus = iota
    UserStatusActive
    UserStatusInactive
    UserStatusBlocked
)
```

## Зависимости

```
cmd                 → internal/app/grpc, internal/modules, internal/adapters, internal/config
internal/app/grpc   → internal/modules (interface), internal/models, internal/types, pkg/pb
internal/modules    → internal/models, internal/types
internal/adapters   → internal/models, internal/types
internal/models     → internal/types
internal/types      → (ничего)
```

**Запрещено:**
- `models` → `modules`
- `modules` → `app/grpc`
- `adapters` → `modules`

## Тестирование

```go
// modules/user_test.go
func TestUserModule_Create(t *testing.T) {
    repo := newMockRepository()
    module := NewUserModule(repo, mockHasher, mockIDGen)

    user, err := module.Create(ctx, models.CreateUserInput{
        Email: "test@example.com",
    })

    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```