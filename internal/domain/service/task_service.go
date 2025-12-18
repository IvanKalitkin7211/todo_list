package service

import (
	"context"
	"github.com/google/uuid"
	"time"
	"todo-list/internal/domain/model"
	"todo-list/internal/domain/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, userID, title, content, status, priority string, due *time.Time) (model.Task, error)
	GetAllTasks(ctx context.Context, userID string) ([]model.Task, error)
	GetTaskByID(ctx context.Context, id, userID string) (model.Task, error)
	UpdateTask(ctx context.Context, id, userID, title, content, status, priority string, due *time.Time) (model.Task, error)
	DeleteTask(ctx context.Context, id, userID string) error
	ChangeStatus(ctx context.Context, id, userID, status string) (model.Task, error)
	GetTasksByStatus(ctx context.Context, status, userID string) ([]model.Task, error)
	SearchTasks(ctx context.Context, q, userID string) ([]model.Task, error)
	GetTodayTasks(ctx context.Context, userID string) ([]model.Task, error)
	GetOverdueTasks(ctx context.Context, userID string) ([]model.Task, error)
	ArchiveTask(ctx context.Context, id, userID string) (model.Task, error)
	UnarchiveTask(ctx context.Context, id, userID string) (model.Task, error)
	ChangePriority(ctx context.Context, id, userID, priority string) (model.Task, error)
	GetTasksByPriority(ctx context.Context, priority, userID string) ([]model.Task, error)
	AddTag(ctx context.Context, id, userID, tag string) (model.Task, error)
	RemoveTag(ctx context.Context, id, userID, tag string) (model.Task, error)
	GetTasksByTag(ctx context.Context, tag, userID string) ([]model.Task, error)
	BulkDelete(ctx context.Context, ids []string, userID string) error
	BulkUpdateStatus(ctx context.Context, ids []string, status, userID string) error
	Stats(ctx context.Context, userID string) (map[string]int64, error)
}

type taskServiceImpl struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskServiceImpl{repo: repo}
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, userID, title, content, status, priority string, due *time.Time) (model.Task, error) {
	uID, _ := uuid.Parse(userID)
	if status == "" {
		status = "todo"
	}
	if priority == "" {
		priority = "medium"
	}
	task := model.Task{
		UserID: uID, Title: title, Content: content, Status: status, Priority: priority, DueDate: due,
	}
	return task, s.repo.Create(ctx, &task)
}

func (s *taskServiceImpl) GetAllTasks(ctx context.Context, userID string) ([]model.Task, error) {
	return s.repo.GetAll(ctx, userID)
}

func (s *taskServiceImpl) GetTaskByID(ctx context.Context, id, userID string) (model.Task, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *taskServiceImpl) UpdateTask(ctx context.Context, id, userID, title, content, status, priority string, due *time.Time) (model.Task, error) {
	task, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return model.Task{}, err
	}
	task.Title = title
	task.Content = content
	task.Status = status
	task.Priority = priority
	task.DueDate = due
	err = s.repo.Update(ctx, &task)
	return task, err
}

func (s *taskServiceImpl) DeleteTask(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *taskServiceImpl) ChangeStatus(ctx context.Context, id, userID, status string) (model.Task, error) {
	task, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return model.Task{}, err
	}
	task.Status = status
	err = s.repo.Update(ctx, &task)
	return task, err
}

func (s *taskServiceImpl) GetTasksByStatus(ctx context.Context, status, userID string) ([]model.Task, error) {
	return s.repo.FindByStatus(ctx, status, userID)
}

func (s *taskServiceImpl) SearchTasks(ctx context.Context, q, userID string) ([]model.Task, error) {
	return s.repo.Search(ctx, q, userID)
}

func (s *taskServiceImpl) GetTodayTasks(ctx context.Context, userID string) ([]model.Task, error) {
	return s.repo.GetToday(ctx, userID)
}

func (s *taskServiceImpl) GetOverdueTasks(ctx context.Context, userID string) ([]model.Task, error) {
	return s.repo.GetOverdue(ctx, userID)
}

func (s *taskServiceImpl) ArchiveTask(ctx context.Context, id, userID string) (model.Task, error) {
	return s.repo.Archive(ctx, id, userID)
}

func (s *taskServiceImpl) UnarchiveTask(ctx context.Context, id, userID string) (model.Task, error) {
	return s.repo.Unarchive(ctx, id, userID)
}

func (s *taskServiceImpl) ChangePriority(ctx context.Context, id, userID, priority string) (model.Task, error) {
	task, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return model.Task{}, err
	}
	task.Priority = priority
	err = s.repo.Update(ctx, &task)
	return task, err
}

func (s *taskServiceImpl) GetTasksByPriority(ctx context.Context, priority, userID string) ([]model.Task, error) {
	return s.repo.FindByPriority(ctx, priority, userID)
}

func (s *taskServiceImpl) AddTag(ctx context.Context, id, userID, tag string) (model.Task, error) {
	return s.repo.AddTag(ctx, id, tag, userID)
}

func (s *taskServiceImpl) RemoveTag(ctx context.Context, id, userID, tag string) (model.Task, error) {
	return s.repo.RemoveTag(ctx, id, tag, userID)
}

func (s *taskServiceImpl) GetTasksByTag(ctx context.Context, tag, userID string) ([]model.Task, error) {
	return s.repo.FindByTag(ctx, tag, userID)
}

func (s *taskServiceImpl) BulkDelete(ctx context.Context, ids []string, userID string) error {
	return s.repo.BulkDelete(ctx, ids, userID)
}

func (s *taskServiceImpl) BulkUpdateStatus(ctx context.Context, ids []string, status, userID string) error {
	return s.repo.BulkUpdateStatus(ctx, ids, status, userID)
}

func (s *taskServiceImpl) Stats(ctx context.Context, userID string) (map[string]int64, error) {
	return s.repo.Stats(ctx, userID)
}
