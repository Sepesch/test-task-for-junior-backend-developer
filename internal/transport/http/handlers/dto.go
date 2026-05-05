package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
}

type taskDTO struct {
	ID           int64             `json:"id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Status       taskdomain.Status `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	RecurrenceID *int64            `json:"recurrence_id,omitempty"`
}

func newTaskDTO(t *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:           t.ID,
		Title:        t.Title,
		Description:  t.Description,
		Status:       t.Status,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		RecurrenceID: t.RecurrenceID,
	}
}

type recurrenceRequest struct {
	Type       string   `json:"type"`
	Interval   *int     `json:"interval,omitempty"`
	DayOfMonth *int     `json:"day_of_month,omitempty"`
	WeekDays   []int    `json:"week_days,omitempty"`
	Dates      []string `json:"dates,omitempty"`
	Parity     *string  `json:"parity,omitempty"`
}

type recurrenceResponse struct {
	ID               int64    `json:"id"`
	Type             string   `json:"type"`
	Interval         *int     `json:"interval,omitempty"`
	DayOfMonth       *int     `json:"day_of_month,omitempty"`
	WeekDays         []int    `json:"week_days,omitempty"`
	Dates            []string `json:"dates,omitempty"`
	Parity           *string  `json:"parity,omitempty"`
	NextOccurrenceAt string   `json:"next_occurrence_at"`
}