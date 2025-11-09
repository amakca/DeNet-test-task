## Структура базы данных (PostgreSQL)

Cхема базы данных сервиса состоит из трёх таблиц:

- **users**: хранит учётные записи пользователей.
- **points**: фиксирует количество баллов пользователя по конкретным заданиям.
- **tasks**: справочник заданий.

## Поля таблиц

- **Таблица users**: `id`, `username`, `password`, `created_at`, `referrer`, `email`
- **Таблица points**: `user_id`, `points`, `task_id`, `upd_at`
- **Таблица tasks**: `id`, `name`, `descr`, `points`

## DDL

```sql
-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
  id          SERIAL PRIMARY KEY,
  username    TEXT        NOT NULL UNIQUE,
  password    TEXT        NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  referrer    INTEGER     NULL REFERENCES users(id) ON DELETE SET NULL,
  email       TEXT UNIQUE
);

-- Таблица заданий
CREATE TABLE IF NOT EXISTS tasks (
  id    SERIAL PRIMARY KEY,
  name  TEXT NOT NULL UNIQUE,
  descr TEXT,
  points INTEGER NOT NULL DEFAULT 0
);

-- Таблица баллов по заданиям
CREATE TABLE IF NOT EXISTS points (
  user_id  INTEGER      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  task_id  INTEGER      NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  points   INTEGER      NOT NULL DEFAULT 0,
  upd_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, task_id)
);

-- Дополнительные индексы
CREATE INDEX IF NOT EXISTS idx_points_user_id ON points(user_id);
CREATE INDEX IF NOT EXISTS idx_points_task_id ON points(task_id);
CREATE INDEX IF NOT EXISTS idx_points_user_id_upd_at ON points(user_id, upd_at DESC);
```

## Связи и ограничения

- `points.user_id` → `users.id` (ON DELETE CASCADE)
- `points.task_id` → `tasks.id` (ON DELETE CASCADE)
- В `points` задан составной первичный ключ `(user_id, task_id)`, чтобы у пользователя была не более одной записи на каждое задание
- Имя задания уникально: в таблице `tasks` добавлено ограничение `UNIQUE (name)`.


