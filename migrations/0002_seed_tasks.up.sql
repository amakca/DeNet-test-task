-- Seed initial tasks with stable IDs used by the service logic
INSERT INTO tasks (id, name, descr, points) VALUES
  (1, 'give_referral',       'User gives referral to a friend',           100),
  (2, 'get_referral',        'User signs up with a referral code',         50),
  (3, 'subscribe_telegram',  'User subscribes to Telegram channel',        30),
  (4, 'subscribe_twitter',   'User follows Twitter/X account',             30),
  (5, 'complete_email',      'User confirms their email address',          20)
ON CONFLICT (id) DO NOTHING;

-- Ensure the sequence is set past the max ID
SELECT setval(pg_get_serial_sequence('tasks', 'id'), (SELECT COALESCE(MAX(id), 0) FROM tasks));


