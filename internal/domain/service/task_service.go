package service

import (
	"context"
	"errors"
	"time"

	"todo-list/internal/domain/model"
	"todo-list/internal/domain/repository"
)

var (
	ErrTaskTitleEmpty   = errors.New("task title is empty")
	ErrTaskContentEmpty = errors.New("task content is empty")
)

type TaskService interface {
	CreateTask(ctx context.Context, title, content, status, priority string, due *time.Time) (model.Task, error)
	GetAllTasks(ctx context.Context) ([]model.Task, error)
	GetTaskByID(ctx context.Context, id string) (model.Task, error)
	UpdateTask(ctx context.Context, id, title, content, status, priority string, due *time.Time) (model.Task, error)
	DeleteTask(ctx context.Context, id string) error

	ChangeStatus(ctx context.Context, id, status string) (model.Task, error)
	GetTasksByStatus(ctx context.Context, status string) ([]model.Task, error)
	SearchTasks(ctx context.Context, q string) ([]model.Task, error)
	GetTodayTasks(ctx context.Context) ([]model.Task, error)
	GetOverdueTasks(ctx context.Context) ([]model.Task, error)
	ArchiveTask(ctx context.Context, id string) (model.Task, error)
	UnarchiveTask(ctx context.Context, id string) (model.Task, error)
	ChangePriority(ctx context.Context, id, priority string) (model.Task, error)
	GetTasksByPriority(ctx context.Context, priority string) ([]model.Task, error)

	AddTag(ctx context.Context, id, tag string) (model.Task, error)
	RemoveTag(ctx context.Context, id, tag string) (model.Task, error)
	GetTasksByTag(ctx context.Context, tag string) ([]model.Task, error)

	BulkDelete(ctx context.Context, ids []string) error
	BulkUpdateStatus(ctx context.Context, ids []string, status string) error

	Stats(ctx context.Context) (map[string]int64, error)
}

type taskServiceImpl struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskServiceImpl{repo: repo}
}

func (s *taskServiceImpl) validateTaskData(title, content string) error {
	if len(title) == 0 {
		return ErrTaskTitleEmpty
	}
	if len(content) == 0 {
		return ErrTaskContentEmpty
	}
	return nil
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, title, content, status, priority string, due *time.Time) (model.Task, error) {
	if err := s.validateTaskData(title, content); err != nil {
		return model.Task{}, err
	}
	if status == "" {
		status = "todo"
	}
	if priority == "" {
		priority = "medium"
	}
	task := model.Task{
		Title:    title,
		Content:  content,
		Status:   status,
		Priority: priority,
		DueDate:  due,
	}
	if err := s.repo.Create(ctx, &task); err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (s *taskServiceImpl) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *taskServiceImpl) GetTaskByID(ctx context.Context, id string) (model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskServiceImpl) UpdateTask(ctx context.Context, id, title, content, status, priority string, due *time.Time) (model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, err
	}
	// validate new values
	if title == "" {
		return model.Task{}, ErrTaskTitleEmpty
	}
	if content == "" {
		return model.Task{}, ErrTaskContentEmpty
	}
	task.Title = title
	task.Content = content
	if status != "" {
		task.Status = status
	}
	if priority != "" {
		task.Priority = priority
	}
	task.DueDate = due
	if err := s.repo.Update(ctx, &task); err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (s *taskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *taskServiceImpl) ChangeStatus(ctx context.Context, id, status string) (model.Task, error) {
	if status == "" {
		return model.Task{}, errors.New("status empty")
	}
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, err
	}
	task.Status = status
	if err := s.repo.Update(ctx, &task); err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (s *taskServiceImpl) GetTasksByStatus(ctx context.Context, status string) ([]model.Task, error) {
	return s.repo.FindByStatus(ctx, status)
}

func (s *taskServiceImpl) SearchTasks(ctx context.Context, q string) ([]model.Task, error) {
	return s.repo.Search(ctx, q)
}

func (s *taskServiceImpl) GetTodayTasks(ctx context.Context) ([]model.Task, error) {
	return s.repo.GetToday(ctx)
}

func (s *taskServiceImpl) GetOverdueTasks(ctx context.Context) ([]model.Task, error) {
	return s.repo.GetOverdue(ctx)
}

func (s *taskServiceImpl) ArchiveTask(ctx context.Context, id string) (model.Task, error) {
	return s.repo.Archive(ctx, id)
}

func (s *taskServiceImpl) UnarchiveTask(ctx context.Context, id string) (model.Task, error) {
	return s.repo.Unarchive(ctx, id)
}

func (s *taskServiceImpl) ChangePriority(ctx context.Context, id, priority string) (model.Task, error) {
	if priority == "" {
		return model.Task{}, errors.New("priority empty")
	}
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, err
	}
	task.Priority = priority
	if err := s.repo.Update(ctx, &task); err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (s *taskServiceImpl) GetTasksByPriority(ctx context.Context, priority string) ([]model.Task, error) {
	return s.repo.FindByPriority(ctx, priority)
}

func (s *taskServiceImpl) AddTag(ctx context.Context, id, tag string) (model.Task, error) {
	return s.repo.AddTag(ctx, id, tag)
}

func (s *taskServiceImpl) RemoveTag(ctx context.Context, id, tag string) (model.Task, error) {
	return s.repo.RemoveTag(ctx, id, tag)
}

func (s *taskServiceImpl) GetTasksByTag(ctx context.Context, tag string) ([]model.Task, error) {
	return s.repo.FindByTag(ctx, tag)
}

func (s *taskServiceImpl) BulkDelete(ctx context.Context, ids []string) error {
	return s.repo.BulkDelete(ctx, ids)
}

func (s *taskServiceImpl) BulkUpdateStatus(ctx context.Context, ids []string, status string) error {
	return s.repo.BulkUpdateStatus(ctx, ids, status)
}

func (s *taskServiceImpl) Stats(ctx context.Context) (map[string]int64, error) {
	return s.repo.Stats(ctx)
}
