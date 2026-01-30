# Архитектурный подход

## Название

**Clean Architecture** с паттерном **Ports & Adapters** и **Vertical Slicing**

## Описание

Слоистая архитектура, где:
- Бизнес-логика изолирована от инфраструктуры через интерфейсы
- Зависимости направлены внутрь (от транспорта к домену)
- Код организован по функциональным единицам (1 файл = 1 use case)

```
┌──────────────────────────────────────────────────────────┐
│                    TRANSPORT LAYER                       │
│              (gRPC handlers, HTTP, CLI)                  │
│         internal/app/grpc/user_service/*.go              │
└──────────────────────────────────────────────────────────┘
                          │
                          ▼
┌──────────────────────────────────────────────────────────┐
│                   BUSINESS LAYER                         │
│              (use cases, бизнес-правила)                 │
│                  internal/modules/                       │
└──────────────────────────────────────────────────────────┘
                          │
            ┌─────────────┴─────────────┐
            ▼                           ▼
┌───────────────────────┐   ┌───────────────────────────────┐
│    DOMAIN LAYER       │   │     INFRASTRUCTURE LAYER      │
│   (бизнес-сущности)   │   │      (БД, API, брокеры)       │
│   internal/models/    │   │      internal/adapters/       │
└───────────────────────┘   └───────────────────────────────┘
            │                           │
            └─────────────┬─────────────┘
                          ▼
              ┌───────────────────────┐
              │    SHARED KERNEL      │
              │   (типы, ошибки)      │
              │   internal/types/     │
              └───────────────────────┘
```

## Ключевые принципы

### 1. Dependency Inversion (Инверсия зависимостей)

Бизнес-логика не зависит от инфраструктуры. Вместо этого:
- `modules/` определяет **интерфейсы** (Ports)
- `adapters/` реализует эти интерфейсы (Adapters)

```go
// internal/modules/user.go — определяет интерфейс
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id string) (*models.User, error)
}

// internal/adapters/repository/postgres.go — реализует
type PostgresUserRepository struct { db *sql.DB }
func (r *PostgresUserRepository) Create(...) error { ... }
```

### 2. Consumer-side Interfaces (Интерфейсы на стороне потребителя)

Интерфейсы определяются там, где используются, а не там, где реализуются.

```
✓ internal/modules/user.go      — type UserRepository interface
✗ internal/adapters/repository/ — НЕ здесь
```

### 3. Vertical Slicing (Вертикальная нарезка)

Один файл = один use case / RPC метод:

```
api/user_service/
├── rpc_create_user.proto      # CreateUser request/response
├── rpc_get_user.proto         # GetUser request/response
└── ...

internal/app/grpc/user_service/
├── create_user.go             # CreateUser handler
├── get_user.go                # GetUser handler
└── ...
```

### 4. Направление зависимостей

```
Transport → Business → Domain
    │           │
    └───────────┴──→ Infrastructure
```

Зависимости всегда направлены **внутрь** (к домену).

## Плюсы для микросервисов

| Преимущество | Описание |
|--------------|----------|
| **Изоляция бизнес-логики** | Смена PostgreSQL на MongoDB не затрагивает `modules/` |
| **Тестируемость** | Unit-тесты бизнес-логики без БД через моки |
| **Независимость от транспорта** | Один `modules/` работает с gRPC, HTTP, Kafka |
| **Чёткие границы** | Каждый слой имеет определённую ответственность |
| **Параллельная разработка** | Команда A делает adapters, команда B — handlers |
| **Простой рефакторинг** | Замена реализации не ломает бизнес-логику |

## Почему удобно для новых разработчиков

### Предсказуемая структура

Любой разработчик сразу знает, где искать:

| Задача | Где искать |
|--------|------------|
| Бизнес-логика | `internal/modules/` |
| Работа с БД | `internal/adapters/repository/` |
| gRPC хендлеры | `internal/app/grpc/<service>/` |
| Модели данных | `internal/models/` |
| Типы ошибок | `internal/types/` |
| Proto файлы | `api/<service>/` |

### Один файл = одна задача

Добавление нового метода `BlockUser`:

```
1. api/user_service/rpc_block_user.proto       — proto
2. internal/app/grpc/user_service/block_user.go — handler
3. internal/modules/user.go                     — бизнес-метод (если нужен)
```

3 файла, каждый с понятной ответственностью.

### Минимальные merge-конфликты

Разработчики работают в разных файлах → меньше конфликтов при слиянии.

### Копипаста как шаблон

Новый метод = скопировать существующий файл и адаптировать:

```bash
cp create_user.go block_user.go
# Заменить CreateUser → BlockUser
```

### Понятный Code Review

PR содержит 3-5 небольших файлов вместо изменений в гигантском `handler.go`.

## Сравнение с другими подходами

| Подход | Плюсы | Минусы |
|--------|-------|--------|
| **MVC** | Простой, знакомый | Бизнес-логика размазана, сложно тестировать |
| **Transaction Script** | Быстро писать | Спагетти-код при росте |
| **Clean Architecture** | Тестируемость, гибкость | Больше файлов, выше порог входа |
| **DDD** | Сложная бизнес-логика | Оверинжиниринг для простых CRUD |

Clean Architecture — баланс между простотой и гибкостью для микросервисов.

## Когда использовать

**Подходит:**
- Микросервисы со сложной бизнес-логикой
- Долгоживущие проекты
- Команды > 2 человек
- Проекты с высокими требованиями к тестируемости

**Избыточно:**
- Простые CRUD-сервисы
- Прототипы / MVP
- Одноразовые скрипты
