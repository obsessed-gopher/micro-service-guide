# Go Microservice Architecture Guide

Демонстрационный проект с архитектурой Go микросервиса.

## Структура проекта

```
.
├── api/                            # Proto файлы
│   └── user_service/
│       ├── api.proto               # Определение сервиса (service UserService)
│       ├── enum.proto              # Enum'ы (UserStatus)
│       ├── model.proto             # Модели (User)
│       ├── rpc_create_user.proto   # CreateUserRequest/Response
│       ├── rpc_get_user.proto      # GetUserRequest/Response
│       ├── rpc_update_user.proto   # UpdateUserRequest/Response
│       ├── rpc_delete_user.proto   # DeleteUserRequest/Response
│       └── rpc_list_users.proto    # ListUsersRequest/Response
│
├── cmd/                            # Точки входа
│   └── user_service/
│       └── main.go
│
├── config/                         # Конфигурации
│   ├── local.yml
│   └── prod.yml
│
├── internal/                       # Внутренний код
│   ├── adapters/                   # Инфраструктура
│   │   ├── hasher/
│   │   ├── idgen/
│   │   └── repository/
│   │
│   ├── config/                     # Структуры конфигурации
│   ├── models/                     # Бизнес-модели
│   ├── modules/                    # Бизнес-логика
│   ├── types/                      # Ошибки, enum'ы
│   ├── metrics/
│   ├── utils/
│   │
│   └── app/                        # Транспортный слой
│       └── grpc/
│           └── user_service/
│               ├── server.go       # Server struct, NewServer()
│               ├── create_user.go  # CreateUser handler
│               ├── get_user.go     # GetUser handler
│               ├── update_user.go  # UpdateUser handler
│               ├── delete_user.go  # DeleteUser handler
│               ├── list_users.go   # ListUsers handler
│               ├── converter.go    # proto ↔ models
│               └── errors.go       # gRPC error mapping
│
├── pkg/pb/                         # Сгенерированный proto код
├── migrations/                     # SQL миграции
├── docker/
└── docs/
```

## Принципы

### 1 файл = 1 сущность

**Proto:**
- `enum.proto` — все enum'ы сервиса
- `model.proto` — модели данных
- `rpc_<method>.proto` — request/response для каждого RPC
- `api.proto` — определение сервиса

**Handlers:**
- `server.go` — структура сервера
- `<method>.go` — один файл на каждый RPC метод
- `converter.go` — конвертеры
- `errors.go` — маппинг ошибок

### Слои

```
cmd/main.go
     │
     ▼
internal/app/grpc/*     # Транспорт (валидация, маппинг)
     │
     ▼
internal/modules/       # Бизнес-логика
     │
     ▼
internal/adapters/      # Инфраструктура (БД, API)
     │
     ▼
internal/models/        # Бизнес-модели
     │
     ▼
internal/types/         # Ошибки, enum'ы
```

## Быстрый старт

```bash
make deps           # Зависимости
make proto          # Генерация proto
make run            # Запуск
make test           # Тесты
make lint           # Линтер
```

## Добавление нового RPC метода

1. Создать `api/user_service/rpc_<method>.proto`
2. Добавить rpc в `api/user_service/api.proto`
3. `make proto`
4. Создать `internal/app/grpc/user_service/<method>.go`

## Добавление нового сервиса

1. Создать `api/new_service/` с proto файлами
2. Создать `internal/app/grpc/new_service/` с хендлерами
3. Добавить модели/модули если нужно
4. Создать `cmd/new_service/main.go`

## Документация

- [docs/ARCHITECTURE_OVERVIEW.md](docs/ARCHITECTURE_OVERVIEW.md) — обзор подхода (Clean Architecture + Ports & Adapters)
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) — детальная архитектура слоёв
- [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) — быстрый старт