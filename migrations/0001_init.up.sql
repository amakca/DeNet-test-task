-- Users table
CREATE TABLE IF NOT EXISTS users (
  id         SERIAL PRIMARY KEY,
  username   TEXT        NOT NULL UNIQUE,
  password   TEXT        NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  referrer   INTEGER     NULL REFERENCES users(id) ON DELETE SET NULL,
  email      TEXT UNIQUE
);

-- Tasks table
CREATE TABLE IF NOT EXISTS tasks (
  id     SERIAL PRIMARY KEY,
  name   TEXT    NOT NULL UNIQUE,
  descr  TEXT,
  points INTEGER NOT NULL DEFAULT 0
);

-- Points table (user progress by tasks)
CREATE TABLE IF NOT EXISTS points (
  user_id INTEGER     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  task_id INTEGER     NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  points  INTEGER     NOT NULL DEFAULT 0,
  upd_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, task_id)
);

-- Helpful indexes
CREATE INDEX IF NOT EXISTS idx_points_user_id ON points(user_id);
CREATE INDEX IF NOT EXISTS idx_points_task_id ON points(task_id);
CREATE INDEX IF NOT EXISTS idx_points_user_id_upd_at ON points(user_id, upd_at DESC);


