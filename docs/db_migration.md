## Миграции базы данных (PostgreSQL)

SQL‑миграции находятся в каталоге `migrations/`.

Использование CLI golang‑migrate:

1) Установка утилиты:

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2) Установите URL вашей базы Postgres (пример для PowerShell):

```powershell
$env:PG_URL="postgres://user:password@localhost:5432/denet?sslmode=disable"
```

3) Применение миграций:

```powershell
migrate -database "$env:PG_URL" -path migrations up
```

4) Откат последней миграции (при необходимости):

```powershell
migrate -database "$env:PG_URL" -path migrations down 1
```

Примечания:
- Создаются таблицы: `users`, `tasks`, `points`.
- Сиды задач (ID 1..5) добавляются в `0002_seed_tasks.up.sql` для соответствия логике сервиса.

### Утилитный скрипт (Windows, PowerShell)

Можно также использовать вспомогательный скрипт:

```powershell
# Задать URL БД на время текущей сессии:
$env:PG_URL="postgres://user:password@localhost:5432/denet?sslmode=disable"

# Применить все up‑миграции
.\scripts\migrate.ps1 -Action up

# Применить N шагов
.\scripts\migrate.ps1 -Action up -Steps 1

# Откатить один шаг
.\scripts\migrate.ps1 -Action down -Steps 1

# Показать текущую версию
.\scripts\migrate.ps1 -Action version

# Форсировать конкретную версию (осторожно)
.\scripts\migrate.ps1 -Action force -Version 2

# Удалить все объекты (осторожно)
.\scripts\migrate.ps1 -Action drop
```

### Автоматические миграции при старте приложения

Приложение применяет миграции при запуске с помощью библиотеки `golang-migrate`:
- Источник: `file://<repo-abs-path>/migrations`
- URL БД берётся из `PG_URL` (см. `config.PG.URL`)

Перед запуском убедитесь, что переменная окружения `PG_URL` задана.