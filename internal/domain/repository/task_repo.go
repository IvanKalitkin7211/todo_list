package repository

import (
	"context"
	"todo-list/internal/domain/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetAll(ctx context.Context, userID string) ([]model.Task, error)
	GetByID(ctx context.Context, id string, userID string) (model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string, userID string) error

	FindByStatus(ctx context.Context, status string, userID string) ([]model.Task, error)
	FindByPriority(ctx context.Context, priority string, userID string) ([]model.Task, error)
	FindByTag(ctx context.Context, tag string, userID string) ([]model.Task, error)
	Search(ctx context.Context, q string, userID string) ([]model.Task, error)
	GetToday(ctx context.Context, userID string) ([]model.Task, error)
	GetOverdue(ctx context.Context, userID string) ([]model.Task, error)

	AddTag(ctx context.Context, id string, tag string, userID string) (model.Task, error)
	RemoveTag(ctx context.Context, id string, tag string, userID string) (model.Task, error)
	BulkDelete(ctx context.Context, ids []string, userID string) error
	BulkUpdateStatus(ctx context.Context, ids []string, status string, userID string) error
	Archive(ctx context.Context, id string, userID string) (model.Task, error)
	Unarchive(ctx context.Context, id string, userID string) (model.Task, error)
	Stats(ctx context.Context, userID string) (map[string]int64, error)
}
