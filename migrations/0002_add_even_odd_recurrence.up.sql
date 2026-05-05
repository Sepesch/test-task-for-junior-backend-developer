-- Таблица правил периодичности
CREATE TABLE task_recurrences (
	id BIGSERIAL PRIMARY KEY,
	type TEXT NOT NULL CHECK (type IN ('daily', 'monthly', 'weekly', 'specific_dates', 'even_odd')),
	config JSONB NOT NULL,
	next_occurrence_at DATE NOT NULL  
);

ALTER TABLE tasks ADD COLUMN recurrence_id BIGINT REFERENCES task_recurrences(id) ON DELETE SET NULL;

CREATE INDEX idx_tasks_recurrence_status ON tasks (recurrence_id, status) WHERE status = 'done';