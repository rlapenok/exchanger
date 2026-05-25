# exchanger

Веб-приложение обменного пункта: курсы валют, операции обмена, журнал действий.

## Роли

| Роль | Возможности |
|------|-------------|
| **Оператор** | Актуальные курсы, «Курсы на дату», обмен валют, журнал своих действий за сессию |
| **Админ** | Ввод и изменение курсов, отчёт по операциям обмена за период |

Тестовые пользователи (создаются миграцией):

| Логин | Пароль | Роль |
|-------|--------|------|
| `admin` | `admin` | admin |
| `operator` | `operator` | operator |

## Быстрый старт (Docker)

Рекомендуемый способ запуска. Миграции применяются автоматически при старте контейнера `app`.

```powershell
docker compose up --build -d
```

Приложение:

```text
http://localhost:8080
```

Наружу открыты только:

| Сервис | URL |
|--------|-----|
| Приложение | http://localhost:8080 |
| pgAdmin | http://localhost:5050 |
| Redis Insight | http://localhost:5540 |

PostgreSQL и Redis доступны только внутри Docker-сети.

Остановка:

```powershell
docker compose stop
```

Полное удаление контейнеров и данных:

```powershell
docker compose down -v
```

## Локальная разработка

### Requirements

- Go `1.25.5`
- Go Task
- `golangci-lint`
- `migrate`
- `templ`

Установка Task:

```powershell
go install github.com/go-task/task/v3/cmd/task@latest
```

Инструменты проекта:

```powershell
task install-tools
```

### Setup

```powershell
task deps
copy .env.example .env
task docker-up
task migrate-up
templ generate ./web/templ/...
go run ./cmd
```

Для локального запуска без Docker-контейнера приложения в `.env` укажите `DB_HOST=localhost`, `REDIS_HOST=localhost` (после `task docker-up`).

Миграции лежат в одном файле: `migrations/000001_init.up.sql`.

## Taskfile

```powershell
task deps           # Download Go module dependencies
task fmt            # Format Go files
task lint           # Run linting checks
task test           # Run tests
task tidy           # Tidy Go modules
task build          # Build the project
task run-binary     # Build and run the binary
task docker-up      # Start Docker Compose services
task docker-stop    # Stop Docker Compose services
task docker-delete  # Remove Docker Compose containers and volumes
task migrate-create NAME=<name> # Create a new SQL migration
task migrate-up     # Apply database migrations
task migrate-down   # Roll back the last database migration
task pre-commit     # Prepare code before commit
task install-tools  # Install development tools
```

## Admin Tools (Docker)

### pgAdmin

```text
http://localhost:5050
```

```text
Email: admin@exchanger.dev
Password: exchanger
```

Подключение к PostgreSQL в pgAdmin:

```text
Host name/address: postgres
Port: 5432
Maintenance database: exchanger
Username: exchanger
Password: exchanger
```

### Redis Insight

```text
http://localhost:5540
```

```text
Host: redis
Port: 6379
Password: exchanger
```

Перед коммитом:

```powershell
task pre-commit
```