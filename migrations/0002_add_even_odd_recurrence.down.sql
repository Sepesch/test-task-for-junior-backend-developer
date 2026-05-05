DROP INDEX IF EXISTS idx_tasks_recurrence_status;
ALTER TABLE tasks DROP COLUMN recurrence_id;
DROP TABLE task_recurrences;