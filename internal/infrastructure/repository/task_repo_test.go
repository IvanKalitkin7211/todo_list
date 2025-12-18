package repository

import (
	"context"
	"testing"
	"time"
	"todo-list/internal/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupRealDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=admin password=normalniy dbname=todolist_db sslmode=disable port=5432"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Чтобы логи базы не спамили в консоль тестов
	})

	if err != nil {
		t.Fatalf("Не удалось подключиться к тестовой БД: %v. Проверь, запущен ли Docker!", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Task{}, &model.Tag{})
	require.NoError(t, err)

	db.Exec("TRUNCATE TABLE task_tags CASCADE")
	db.Exec("TRUNCATE TABLE tasks CASCADE")
	db.Exec("TRUNCATE TABLE users CASCADE")

	return db
}

func TestRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест в коротком режиме")
	}

	db := setupRealDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	db.Exec("INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)", userID, "test@example.com", "hash")

	t.Run("Full_Task_Lifecycle", func(t *testing.T) {
		taskID := uuid.New()
		task := &model.Task{
			ID:       taskID,
			UserID:   userID,
			Title:    "Real DB Task",
			Content:  "Integration testing",
			Status:   "todo",
			Priority: "high",
		}

		// 1. Create
		err := repo.Create(ctx, task)
		assert.NoError(t, err)

		// 2. Get and Verify
		savedTask, err := repo.GetByID(ctx, taskID.String(), userID.String())
		assert.NoError(t, err)
		assert.Equal(t, "Real DB Task", savedTask.Title)

		// 3. Stats
		stats, err := repo.Stats(ctx, userID.String())
		assert.NoError(t, err)
		assert.Equal(t, int64(1), stats["todo"])
	})

	t.Run("Bulk_Delete_Integration", func(t *testing.T) {
		id1, id2 := uuid.New(), uuid.New()
		db.Create(&model.Task{ID: id1, UserID: userID, Title: "T1"})
		db.Create(&model.Task{ID: id2, UserID: userID, Title: "T2"})

		err := repo.BulkDelete(ctx, []string{id1.String(), id2.String()}, userID.String())
		assert.NoError(t, err)

		var count int64
		db.Model(&model.Task{}).Where("id IN ?", []uuid.UUID{id1, id2}).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestRepository_GetByID(t *testing.T) {
	db := setupRealDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()

	uid := uuid.New()
	taskID := uuid.New()
	originalTask := &model.Task{
		ID:      taskID,
		UserID:  uid,
		Title:   "Test Task",
		Content: "Content",
	}
	repo.Create(ctx, originalTask)

	t.Run("Success", func(t *testing.T) {
		task, err := repo.GetByID(ctx, taskID.String(), uid.String())

		assert.NoError(t, err)
		assert.Equal(t, "Test Task", task.Title)
		assert.Equal(t, taskID, task.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := repo.GetByID(ctx, uuid.New().String(), uid.String())

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRepository_Features(t *testing.T) {
	db := setupRealDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()
	uid := uuid.New()
	userID := uid.String()

	t.Run("Search", func(t *testing.T) {
		repo.Create(ctx, &model.Task{ID: uuid.New(), UserID: uid, Title: "Купить молоко"})

		results, err := repo.Search(ctx, "МОЛОКО", userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Tags", func(t *testing.T) {
		task := &model.Task{ID: uuid.New(), UserID: uid, Title: "Task with Tag"}
		repo.Create(ctx, task)

		updated, err := repo.AddTag(ctx, task.ID.String(), "urgent", userID)
		assert.NoError(t, err)
		assert.Len(t, updated.Tags, 1)

		tasksByTag, err := repo.FindByTag(ctx, "urgent", userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, tasksByTag)
	})

	// 3. Тестируем Статистику (Count)
	t.Run("Stats", func(t *testing.T) {
		stats, err := repo.Stats(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Contains(t, stats, "total")
	})
}

func TestRepository_TimeQueries(t *testing.T) {
	db := setupRealDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()
	uid := uuid.New()

	t.Run("GetToday_And_Overdue", func(t *testing.T) {
		now := time.Now()
		repo.Create(ctx, &model.Task{ID: uuid.New(), UserID: uid, Title: "Today", DueDate: &now, Status: "todo"})

		yesterday := now.AddDate(0, 0, -1)
		repo.Create(ctx, &model.Task{ID: uuid.New(), UserID: uid, Title: "Old", DueDate: &yesterday, Status: "todo"})

		todayTasks, err := repo.GetToday(ctx, uid.String())
		assert.NoError(t, err)
		assert.Len(t, todayTasks, 1)

		overdueTasks, err := repo.GetOverdue(ctx, uid.String())
		assert.NoError(t, err)
		assert.NotEmpty(t, overdueTasks)
	})
}

func TestRepository_BulkOperations(t *testing.T) {
	db := setupRealDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()
	uid := uuid.New()

	id1, id2 := uuid.New(), uuid.New()
	repo.Create(ctx, &model.Task{ID: id1, UserID: uid, Status: "todo"})
	repo.Create(ctx, &model.Task{ID: id2, UserID: uid, Status: "todo"})

	err := repo.BulkUpdateStatus(ctx, []string{id1.String(), id2.String()}, "done", uid.String())
	assert.NoError(t, err)

	task, _ := repo.GetByID(ctx, id1.String(), uid.String())
	assert.Equal(t, "done", task.Status)
}

func TestRepository_Archive(t *testing.T) {
	db := setupRealDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()
	uid := uuid.New()
	tid := uuid.New()

	repo.Create(ctx, &model.Task{ID: tid, UserID: uid, Archived: false})

	// Archive
	res, err := repo.Archive(ctx, tid.String(), uid.String())
	assert.NoError(t, err)
	assert.True(t, res.Archived)

	// Unarchive
	res, err = repo.Unarchive(ctx, tid.String(), uid.String())
	assert.NoError(t, err)
	assert.False(t, res.Archived)
}
