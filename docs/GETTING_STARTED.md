# Быстрый старт

## Требования

- Go 1.22+
- Docker & Docker Compose
- protoc + protoc-gen-go + protoc-gen-go-grpc
- golangci-lint

```bash
# macOS
brew install go protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Запуск

```bash
# Зависимости
make deps

# Генерация proto
make proto

# Запуск
make run

# Тесты
make test
```

## Добавление нового RPC метода

### 1. Создать proto файл

```bash
touch api/user_service/rpc_block_user.proto
```

```protobuf
// api/user_service/rpc_block_user.proto
syntax = "proto3";
package user_service;
option go_package = "github.com/example/user-service/pkg/pb/user_service";

message BlockUserRequest {
  string id = 1;
  string reason = 2;
}

message BlockUserResponse {}
```

### 2. Добавить в api.proto

```protobuf
// api/user_service/api.proto
import "api/user_service/rpc_block_user.proto";

service UserService {
  // ... existing rpcs
  rpc BlockUser(BlockUserRequest) returns (BlockUserResponse);
}
```

### 3. Сгенерировать код

```bash
make proto
```

### 4. Добавить метод в modules (если нужна новая логика)

```go
// internal/modules/user.go
func (m *UserModule) Block(ctx context.Context, id, reason string) error {
    user, err := m.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    user.Status = types.UserStatusBlocked
    user.UpdatedAt = time.Now()

    return m.repo.Update(ctx, user)
}
```

### 5. Создать handler

```bash
touch internal/app/grpc/user_service/block_user.go
```

```go
// internal/app/grpc/user_service/block_user.go
package user_service

import (
    "context"

    pb "github.com/example/user-service/pkg/pb/user_service"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *Server) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "id is required")
    }

    if err := s.userModule.Block(ctx, req.Id, req.Reason); err != nil {
        return nil, mapError(err)
    }

    return &pb.BlockUserResponse{}, nil
}
```

### 6. Обновить интерфейс в server.go (если добавили метод в modules)

```go
// internal/app/grpc/user_service/server.go
type UserModule interface {
    // ... existing methods
    Block(ctx context.Context, id, reason string) error
}
```

## Добавление нового сервиса

### 1. Создать proto файлы

```bash
mkdir -p api/order_service
```

```
api/order_service/
├── api.proto
├── enum.proto
├── model.proto
├── rpc_create_order.proto
├── rpc_get_order.proto
└── rpc_list_orders.proto
```

### 2. Создать handlers

```bash
mkdir -p internal/app/grpc/order_service
```

```
internal/app/grpc/order_service/
├── server.go
├── create_order.go
├── get_order.go
├── list_orders.go
├── converter.go
└── errors.go
```

### 3. Добавить модели и бизнес-логику

```
internal/models/order.go
internal/modules/order.go
internal/adapters/repository/order_postgres.go
```

### 4. Создать точку входа

```bash
mkdir -p cmd/order_service
touch cmd/order_service/main.go
```

### 5. Обновить Makefile

```makefile
.PHONY: run-order
run-order:
	go run ./cmd/order_service -config=config/local.yml
```

## Структура файлов

### Proto (1 RPC = 1 файл)

```
api/user_service/
├── api.proto               # service определение
├── enum.proto              # все enum'ы
├── model.proto             # общие модели (User, etc.)
├── rpc_create_user.proto   # CreateUserRequest/Response
├── rpc_get_user.proto
├── rpc_update_user.proto
├── rpc_delete_user.proto
└── rpc_list_users.proto
```

### Handlers (1 RPC = 1 файл)

```
internal/app/grpc/user_service/
├── server.go               # Server struct, NewServer()
├── create_user.go          # func (s *Server) CreateUser()
├── get_user.go
├── update_user.go
├── delete_user.go
├── list_users.go
├── converter.go            # userToProto(), statusFromProto()
└── errors.go               # mapError()
```

## Полезные команды

```bash
# Запуск с hot-reload (нужен air)
air

# Тесты с покрытием
make test-coverage

# gRPC клиент
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"email":"test@test.com","name":"Test","password":"12345678"}' \
  localhost:50051 user_service.UserService/CreateUser

# Docker
docker-compose -f docker/docker-compose.yml up
```