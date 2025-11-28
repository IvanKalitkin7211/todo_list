package dto

import (
	"time"
	"todo-list/internal/domain/model"
)

type TaskResponseDTO struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Status    string     `json:"status"`
	Priority  string     `json:"priority"`
	Tags      []string   `json:"tags,omitempty"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	Archived  bool       `json:"archived"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func ToTaskResponseDTO(task model.Task) TaskResponseDTO {
	var tags []string
	for _, t := range task.Tags {
		tags = append(tags, t.Name)
	}
	return TaskResponseDTO{
		ID:        task.ID.String(),
		Title:     task.Title,
		Content:   task.Content,
		Status:    task.Status,
		Priority:  task.Priority,
		Tags:      tags,
		DueDate:   task.DueDate,
		Archived:  task.Archived,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
