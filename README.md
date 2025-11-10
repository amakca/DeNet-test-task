# DeNet-test-task

### Описание
Сервис начисляет баллы пользователям за выполнение заданий (подписки, рефералы и т.п.). Есть:
- аутентификация (JWT),
- сущности пользователей, заданий и баллов,
- рейтинги (лидерборд),
- история начислений,
- миграции базы данных на PostgreSQL.

### Стек
- Go (минимальная версия — из `go.mod`)
- PostgreSQL
- Chi (HTTP‑роутер)
- pgx + Squirrel (доступ к БД)
- golang-migrate (миграции)
- slog (логирование)

### Структура проекта
- `cmd/app` — точка входа (main)
- `config` — загрузка конфигурации
- `internal/app` — инициализация приложения (логирование, БД, миграции, HTTP‑сервер)
- `internal/api/v1` — HTTP‑роуты, middleware, хендлеры
- `internal/services` — бизнес‑логика (auth, users, tasks)
- `internal/repo` — интерфейсы и реализации репозиториев (`internal/repo/pgdb`)
- `internal/entity` — доменные структуры (`User`, `Task`, `Point`)
- `pkg/postgres` — обёртка над pgx и билдером запросов
- `pkg/httpserver` — HTTP‑сервер
- `pkg/hasher` — хеширование паролей (с солью)
- `pkg/validator` — простая валидация
- `pkg/migrator` — программный раннер миграций (golang‑migrate)
- `migrations/` — SQL‑миграции
- `docs/` — документация (`db_schema.md`, `db_migration.md`)

### Требования
- Установленный Go
- Доступный PostgreSQL (локально или в контейнере)

### Конфигурация
Базовый файл: `cmd/app/config/config.yaml`

Обязательные переменные окружения:
- `PG_URL` — строка подключения к БД PostgreSQL
- `JWT_SIGN_KEY` — секрет для подписи JWT
- `HASHER_SALT` — соль для хеширования паролей

Пример для PowerShell (Windows):
```powershell
$env:PG_URL="postgres://user:password@localhost:5432/denet?sslmode=disable"
$env:JWT_SIGN_KEY="super-secret"
$env:HASHER_SALT="my-salt"
```

Логи:
- Человекочитаемые (text) при `ENV=dev|development`
- JSON по умолчанию

### Миграции БД
- Миграции находятся в `migrations/`
- На старте приложения миграции применяются автоматически (golang‑migrate)
- Дополнительно можно запускать вручную: см. `docs/db_migration.md`

Сиды заданий находятся в `0002_seed_tasks.up.sql` (ID 1..5 фиксированы и используются сервисом).

### Сборка и запуск
Запуск:
```bash
go run ./cmd/app
```

Сборка:
```bash
go build -o bin/app ./cmd/app
./bin/app
```

Перед запуском убедитесь, что заданы переменные окружения (`PG_URL`, `JWT_SIGN_KEY`, `HASHER_SALT`) и доступна БД.

### HTTP API (кратко)
- `GET /health` — проверка живости

Аутентификация:
- `POST /auth/sign-up` — регистрация пользователя
  - тело: `{ "username": "u", "password": "p" }`
  - ответ: `{ "id": 1 }`
- `POST /auth/sign-in` — получение JWT
  - тело: `{ "username": "u", "password": "p" }`
  - ответ: `{ "token": "<jwt>" }`

Все эндпоинты ниже требуют заголовок `Authorization: Bearer <jwt>`.

Пользователи (`/api/v1/users`):
- `GET /{user_id}/status` — информация о пользователе
- `GET /{user_id}/history?limit=N` — история начислений (1..100)
- `GET /{user_id}/points` — суммарные баллы
- `GET /leaderboard?limit=N` — лидерборд
- `POST /{user_id}/referrer` — задать реферера (form: `referrer=<id>`)
- `POST /{user_id}/email` — задать email (form: `email=<value>`)
- `POST /{user_id}/task/complete` — завершить задание (form: `task_id=<id>`)

Задания (`/api/v1/tasks`):
- `GET /list` — список заданий

### Схема БД
Краткое описание таблиц и связей: см. `docs/db_schema.md`.

### Разработка
- Логи уровня debug: `ENV=dev`
- Ручные миграции и утилиты: `docs/db_migration.md`, `scripts/migrate.ps1`

### Лицензия
Не указана (по необходимости добавьте).
