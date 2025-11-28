package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"todo-list/internal/domain/model"
	drepo "todo-list/internal/domain/repository"

	"gorm.io/gorm"
)

type taskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) drepo.TaskRepository {
	// return typed interface
	return &taskRepositoryImpl{db: db}
}

func (r *taskRepositoryImpl) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepositoryImpl) GetAll(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepositoryImpl) GetByID(ctx context.Context, id string) (model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").First(&task, "id = ?", id).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *taskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(task).Error
}

func (r *taskRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Task{}, "id = ?", id).Error
}

func (r *taskRepositoryImpl) FindByStatus(ctx context.Context, status string) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").Where("status = ?", status).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepositoryImpl) FindByPriority(ctx context.Context, priority string) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").Where("priority = ?", priority).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepositoryImpl) FindByTag(ctx context.Context, tag string) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.WithContext(ctx).
		Joins("JOIN task_tags ON task_tags.task_id = tasks.id").
		Joins("JOIN tags ON tags.id = task_tags.tag_id AND tags.name = ?", tag).
		Preload("Tags").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepositoryImpl) Search(ctx context.Context, q string) ([]model.Task, error) {
	var tasks []model.Task
	q = strings.TrimSpace(q)
	if q == "" {
		return []model.Task{}, nil
	}
	like := fmt.Sprintf("%%%s%%", q)
	if err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("title ILIKE ? OR content ILIKE ?", like, like).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepositoryImpl) GetToday(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24*time.Hour - time.Nanosecond)
	if err := r.db.WithContext(ctx).Preload("Tags").Where("due_date >= ? AND due_date <= ?", start, end).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepositoryImpl) GetOverdue(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").Where("due_date < ? AND status <> ?", time.Now(), "done").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepositoryImpl) AddTag(ctx context.Context, id string, tag string) (model.Task, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return model.Task{}, errors.New("tag empty")
	}
	var task model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").First(&task, "id = ?", id).Error; err != nil {
		return model.Task{}, err
	}
	var t model.Tag
	if err := r.db.WithContext(ctx).FirstOrCreate(&t, model.Tag{Name: tag}).Error; err != nil {
		return model.Task{}, err
	}
	if err := r.db.WithContext(ctx).Model(&task).Association("Tags").Append(&t); err != nil {
		return model.Task{}, err
	}
	if err := r.db.WithContext(ctx).Preload("Tags").First(&task, "id = ?", id).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *taskRepositoryImpl) RemoveTag(ctx context.Context, id string, tag string) (model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").First(&task, "id = ?", id).Error; err != nil {
		return model.Task{}, err
	}
	var t model.Tag
	if err := r.db.WithContext(ctx).First(&t, "name = ?", tag).Error; err != nil {
		return model.Task{}, err
	}
	if err := r.db.WithContext(ctx).Model(&task).Association("Tags").Delete(&t); err != nil {
		return model.Task{}, err
	}
	if err := r.db.WithContext(ctx).Preload("Tags").First(&task, "id = ?", id).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *taskRepositoryImpl) BulkDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Task{}).Error
}

func (r *taskRepositoryImpl) BulkUpdateStatus(ctx context.Context, ids []string, status string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&model.Task{}).Where("id IN ?", ids).Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *taskRepositoryImpl) Archive(ctx context.Context, id string) (model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").First(&task, "id = ?", id).Error; err != nil {
		return model.Task{}, err
	}
	task.Archived = true
	if err := r.db.WithContext(ctx).Save(&task).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *taskRepositoryImpl) Unarchive(ctx context.Context, id string) (model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).Preload("Tags").First(&task, "id = ?", id).Error; err != nil {
		return model.Task{}, err
	}
	task.Archived = false
	if err := r.db.WithContext(ctx).Save(&task).Error; err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *taskRepositoryImpl) Stats(ctx context.Context) (map[string]int64, error) {
	out := map[string]int64{}
	statuses := []string{"todo", "in_progress", "done", "blocked"}
	var cnt int64
	for _, s := range statuses {
		if err := r.db.WithContext(ctx).Model(&model.Task{}).Where("status = ?", s).Count(&cnt).Error; err != nil {
			return nil, err
		}
		out[s] = cnt
	}
	// total
	if err := r.db.WithContext(ctx).Model(&model.Task{}).Count(&cnt).Error; err != nil {
		return nil, err
	}
	out["total"] = cnt
	return out, nil
}
