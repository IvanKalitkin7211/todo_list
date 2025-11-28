package repository

import (
	"context"
	"todo-list/internal/domain/model"
)

// TaskRepository — интерфейс для репозитория задач
type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetAll(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id string) (model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error

	FindByStatus(ctx context.Context, status string) ([]model.Task, error)
	FindByPriority(ctx context.Context, priority string) ([]model.Task, error)
	FindByTag(ctx context.Context, tag string) ([]model.Task, error)
	Search(ctx context.Context, q string) ([]model.Task, error)
	GetToday(ctx context.Context) ([]model.Task, error)
	GetOverdue(ctx context.Context) ([]model.Task, error)

	AddTag(ctx context.Context, id string, tag string) (model.Task, error)
	RemoveTag(ctx context.Context, id string, tag string) (model.Task, error)
	BulkDelete(ctx context.Context, ids []string) error
	BulkUpdateStatus(ctx context.Context, ids []string, status string) error
	Archive(ctx context.Context, id string) (model.Task, error)
	Unarchive(ctx context.Context, id string) (model.Task, error)
	Stats(ctx context.Context) (map[string]int64, error)
}
