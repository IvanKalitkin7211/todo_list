package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"todo-list/internal/domain/model"
	drepo "todo-list/internal/domain/repository"
)

type taskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) drepo.TaskRepository {
	return &taskRepositoryImpl{db: db}
}

func (r *taskRepositoryImpl) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepositoryImpl) GetAll(ctx context.Context, userID string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).Preload("Tags").Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepositoryImpl) GetByID(ctx context.Context, id string, userID string) (model.Task, error) {
	var task model.Task
	err := r.db.WithContext(ctx).Preload("Tags").Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	return task, err
}

func (r *taskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(task).Error
}

func (r *taskRepositoryImpl) Delete(ctx context.Context, id string, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Task{}).Error
}

func (r *taskRepositoryImpl) FindByStatus(ctx context.Context, status string, userID string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).Preload("Tags").Where("status = ? AND user_id = ?", status, userID).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepositoryImpl) FindByPriority(ctx context.Context, priority string, userID string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).Preload("Tags").Where("priority = ? AND user_id = ?", priority, userID).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepositoryImpl) FindByTag(ctx context.Context, tag string, userID string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).
		Joins("JOIN task_tags ON task_tags.task_id = tasks.id").
		Joins("JOIN tags ON tags.id = task_tags.tag_id AND tags.name = ?", tag).
		Preload("Tags").
		Where("tasks.user_id = ?", userID).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepositoryImpl) Search(ctx context.Context, q string, userID string) ([]model.Task, error) {
	var tasks []model.Task
	like := fmt.Sprintf("%%%s%%", strings.TrimSpace(q))
	err := r.db.WithContext(ctx).Preload("Tags").
		Where("user_id = ? AND (title ILIKE ? OR content ILIKE ?)", userID, like, like).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepositoryImpl) GetToday(ctx context.Context, userID string) ([]model.Task, error) {
	var tasks []model.Task
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24*time.Hour - time.Nanosecond)
	err := r.db.WithContext(ctx).Preload("Tags").
		Where("user_id = ? AND due_date >= ? AND due_date <= ?", userID, start, end).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepositoryImpl) GetOverdue(ctx context.Context, userID string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).Preload("Tags").
		Where("user_id = ? AND due_date < ? AND status <> ?", userID, time.Now(), "done").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepositoryImpl) AddTag(ctx context.Context, id string, tag string, userID string) (model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		return model.Task{}, err
	}
	var t model.Tag
	r.db.WithContext(ctx).FirstOrCreate(&t, model.Tag{Name: tag})
	r.db.WithContext(ctx).Model(&task).Association("Tags").Append(&t)
	return r.GetByID(ctx, id, userID)
}

func (r *taskRepositoryImpl) RemoveTag(ctx context.Context, id string, tag string, userID string) (model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		return model.Task{}, err
	}
	var t model.Tag
	if err := r.db.WithContext(ctx).Where("name = ?", tag).First(&t).Error; err == nil {
		r.db.WithContext(ctx).Model(&task).Association("Tags").Delete(&t)
	}
	return r.GetByID(ctx, id, userID)
}

func (r *taskRepositoryImpl) BulkDelete(ctx context.Context, ids []string, userID string) error {
	return r.db.WithContext(ctx).Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.Task{}).Error
}

func (r *taskRepositoryImpl) BulkUpdateStatus(ctx context.Context, ids []string, status string, userID string) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).
		Where("id IN ? AND user_id = ?", ids, userID).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *taskRepositoryImpl) Archive(ctx context.Context, id string, userID string) (model.Task, error) {
	err := r.db.WithContext(ctx).Model(&model.Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("archived", true).Error
	if err != nil {
		return model.Task{}, err
	}
	return r.GetByID(ctx, id, userID)
}

func (r *taskRepositoryImpl) Unarchive(ctx context.Context, id string, userID string) (model.Task, error) {
	err := r.db.WithContext(ctx).Model(&model.Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("archived", false).Error
	if err != nil {
		return model.Task{}, err
	}
	return r.GetByID(ctx, id, userID)
}

func (r *taskRepositoryImpl) Stats(ctx context.Context, userID string) (map[string]int64, error) {
	out := map[string]int64{}
	statuses := []string{"todo", "in_progress", "done", "blocked"}
	for _, s := range statuses {
		var cnt int64
		r.db.WithContext(ctx).Model(&model.Task{}).Where("status = ? AND user_id = ?", s, userID).Count(&cnt)
		out[s] = cnt
	}
	var total int64
	r.db.WithContext(ctx).Model(&model.Task{}).Where("user_id = ?", userID).Count(&total)
	out["total"] = total
	return out, nil
}
