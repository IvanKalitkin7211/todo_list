package testutils

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
	"todo-list/internal/domain/model"
)

type AllMocks struct {
	mock.Mock
}

// Репозиторий
func (m *AllMocks) Create(ctx context.Context, t *model.Task) error { return m.Called(ctx, t).Error(0) }
func (m *AllMocks) GetAll(ctx context.Context, uID string) ([]model.Task, error) {
	args := m.Called(ctx, uID)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) GetByID(ctx context.Context, id, uID string) (model.Task, error) {
	args := m.Called(ctx, id, uID)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) Update(ctx context.Context, t *model.Task) error { return m.Called(ctx, t).Error(0) }
func (m *AllMocks) Delete(ctx context.Context, id, uID string) error {
	return m.Called(ctx, id, uID).Error(0)
}
func (m *AllMocks) FindByStatus(ctx context.Context, s, uID string) ([]model.Task, error) {
	args := m.Called(ctx, s, uID)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) FindByPriority(ctx context.Context, p, uID string) ([]model.Task, error) {
	args := m.Called(ctx, p, uID)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) FindByTag(ctx context.Context, t, uID string) ([]model.Task, error) {
	args := m.Called(ctx, t, uID)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) Search(ctx context.Context, q, uID string) ([]model.Task, error) {
	args := m.Called(ctx, q, uID)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) GetToday(ctx context.Context, uID string) ([]model.Task, error) {
	args := m.Called(ctx, uID)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) GetOverdue(ctx context.Context, uID string) ([]model.Task, error) {
	args := m.Called(ctx, uID)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) AddTag(ctx context.Context, id, tag, uID string) (model.Task, error) {
	args := m.Called(ctx, id, tag, uID)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) RemoveTag(ctx context.Context, id, tag, uID string) (model.Task, error) {
	args := m.Called(ctx, id, tag, uID)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) BulkDelete(ctx context.Context, ids []string, uID string) error {
	return m.Called(ctx, ids, uID).Error(0)
}
func (m *AllMocks) BulkUpdateStatus(ctx context.Context, ids []string, s, uID string) error {
	return m.Called(ctx, ids, s, uID).Error(0)
}
func (m *AllMocks) Archive(ctx context.Context, id, uID string) (model.Task, error) {
	args := m.Called(ctx, id, uID)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) Unarchive(ctx context.Context, id, uID string) (model.Task, error) {
	args := m.Called(ctx, id, uID)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) Stats(ctx context.Context, uID string) (map[string]int64, error) {
	args := m.Called(ctx, uID)
	return args.Get(0).(map[string]int64), args.Error(1)
}

// Сервис (методы CreateTask и т.д.)
func (m *AllMocks) CreateTask(ctx context.Context, u, t, c, s, p string, d *time.Time) (model.Task, error) {
	args := m.Called(ctx, u, t, c, s, p, d)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) GetAllTasks(ctx context.Context, u string) ([]model.Task, error) {
	args := m.Called(ctx, u)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) GetTaskByID(ctx context.Context, id, u string) (model.Task, error) {
	args := m.Called(ctx, id, u)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) UpdateTask(ctx context.Context, id, u, t, c, s, p string, d *time.Time) (model.Task, error) {
	args := m.Called(ctx, id, u, t, c, s, p, d)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) DeleteTask(ctx context.Context, id, u string) error {
	return m.Called(ctx, id, u).Error(0)
}
func (m *AllMocks) ChangeStatus(ctx context.Context, id, u, s string) (model.Task, error) {
	args := m.Called(ctx, id, u, s)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) GetTasksByStatus(ctx context.Context, s, u string) ([]model.Task, error) {
	args := m.Called(ctx, s, u)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) SearchTasks(ctx context.Context, q, u string) ([]model.Task, error) {
	args := m.Called(ctx, q, u)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) GetTodayTasks(ctx context.Context, u string) ([]model.Task, error) {
	args := m.Called(ctx, u)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) GetOverdueTasks(ctx context.Context, u string) ([]model.Task, error) {
	args := m.Called(ctx, u)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) ArchiveTask(ctx context.Context, id, u string) (model.Task, error) {
	args := m.Called(ctx, id, u)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) UnarchiveTask(ctx context.Context, id, u string) (model.Task, error) {
	args := m.Called(ctx, id, u)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) ChangePriority(ctx context.Context, id, u, p string) (model.Task, error) {
	args := m.Called(ctx, id, u, p)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) GetTasksByPriority(ctx context.Context, p, u string) ([]model.Task, error) {
	args := m.Called(ctx, p, u)
	return args.Get(0).([]model.Task), args.Error(1)
}
func (m *AllMocks) AddTagToTask(ctx context.Context, id, u, t string) (model.Task, error) {
	args := m.Called(ctx, id, u, t)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) RemoveTagFromTask(ctx context.Context, id, u, t string) (model.Task, error) {
	args := m.Called(ctx, id, u, t)
	return args.Get(0).(model.Task), args.Error(1)
}
func (m *AllMocks) GetTasksByTag(ctx context.Context, t, u string) ([]model.Task, error) {
	args := m.Called(ctx, t, u)
	return args.Get(0).([]model.Task), args.Error(1)
}
