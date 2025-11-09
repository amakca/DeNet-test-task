-- Remove seeded tasks by IDs
DELETE FROM tasks WHERE id IN (1, 2, 3, 4, 5);

-- Reseat the sequence to the current max
SELECT setval(pg_get_serial_sequence('tasks', 'id'), (SELECT COALESCE(MAX(id), 0) FROM tasks));


