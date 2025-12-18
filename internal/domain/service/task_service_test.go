package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"todo-list/internal/domain/model"
	"todo-list/internal/testutils"
)

func TestTaskService_FullSuite(t *testing.T) {
	repo := new(testutils.AllMocks)
	svc := NewTaskService(repo)
	ctx := context.Background()
	uID := uuid.New().String()

	t.Run("CreateTask_Valid", func(t *testing.T) {
		repo.On("Create", ctx, mock.AnythingOfType("*model.Task")).Return(nil).Once()
		res, err := svc.CreateTask(ctx, uID, "Title", "Content", "todo", "high", nil)
		assert.NoError(t, err)
		assert.Equal(t, "Title", res.Title)
	})

	t.Run("CreateTask_DefaultValues", func(t *testing.T) {
		repo.On("Create", ctx, mock.MatchedBy(func(task *model.Task) bool {
			return task.Status == "todo" && task.Priority == "medium"
		})).Return(nil).Once()
		res, err := svc.CreateTask(ctx, uID, "T", "C", "", "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "todo", res.Status)
	})

	t.Run("GetTaskByID_Success", func(t *testing.T) {
		tID := uuid.New().String()
		repo.On("GetByID", ctx, tID, uID).Return(model.Task{Title: "X"}, nil).Once()
		res, err := svc.GetTaskByID(ctx, tID, uID)
		assert.NoError(t, err)
		assert.Equal(t, "X", res.Title)
	})

	t.Run("GetTaskByID_NotFound", func(t *testing.T) {
		repo.On("GetByID", ctx, "invalid", uID).Return(model.Task{}, errors.New("not found")).Once()
		_, err := svc.GetTaskByID(ctx, "invalid", uID)
		assert.Error(t, err)
	})

	t.Run("Stats_Calculation", func(t *testing.T) {
		mockStats := map[string]int64{"todo": 2, "done": 5}
		repo.On("Stats", ctx, uID).Return(mockStats, nil).Once()
		res, err := svc.Stats(ctx, uID)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), res["done"])
	})

	t.Run("BulkDelete_Execute", func(t *testing.T) {
		ids := []string{"1", "2"}
		repo.On("BulkDelete", ctx, ids, uID).Return(nil).Once()
		err := svc.BulkDelete(ctx, ids, uID)
		assert.NoError(t, err)
	})

	t.Run("SearchTasks_Success", func(t *testing.T) {
		query := "milk"
		repo.On("Search", ctx, query, uID).Return([]model.Task{{Title: "Buy milk"}}, nil).Once()
		res, err := svc.SearchTasks(ctx, query, uID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "Buy milk", res[0].Title)
	})

	t.Run("Tag_Operations", func(t *testing.T) {
		tID := uuid.New().String()
		uID := uuid.New().String()
		tagName := "urgent"

		repo.On("AddTag", ctx, tID, tagName, uID).Return(model.Task{Title: "T"}, nil).Once()

		_, err := svc.AddTag(ctx, tID, uID, tagName)
		assert.NoError(t, err)

		repo.On("RemoveTag", ctx, tID, tagName, uID).Return(model.Task{Title: "T"}, nil).Once()
		_, err = svc.RemoveTag(ctx, tID, uID, tagName)
		assert.NoError(t, err)
	})

	t.Run("Archive_Lifecycle", func(t *testing.T) {
		tID := uuid.New().String()

		// Archive
		repo.On("Archive", ctx, tID, uID).Return(model.Task{Archived: true}, nil).Once()
		res, err := svc.ArchiveTask(ctx, tID, uID)
		assert.NoError(t, err)
		assert.True(t, res.Archived)

		// Unarchive
		repo.On("Unarchive", ctx, tID, uID).Return(model.Task{Archived: false}, nil).Once()
		res, err = svc.UnarchiveTask(ctx, tID, uID)
		assert.NoError(t, err)
		assert.False(t, res.Archived)
	})

	t.Run("Priority_And_Status_Changes", func(t *testing.T) {
		tID := uuid.New().String()
		uID := uuid.New().String()

		// ChangePriority
		repo.On("GetByID", ctx, tID, uID).Return(model.Task{ID: [16]byte{1}}, nil).Once()
		repo.On("Update", ctx, mock.Anything).Return(nil).Once()
		_, err := svc.ChangePriority(ctx, tID, uID, "high")
		assert.NoError(t, err)

		// BulkUpdateStatus
		ids := []string{tID}
		repo.On("BulkUpdateStatus", ctx, ids, "done", uID).Return(nil).Once()
		err = svc.BulkUpdateStatus(ctx, ids, "done", uID)
		assert.NoError(t, err)
	})

	t.Run("List_Filters", func(t *testing.T) {
		uID := uuid.New().String()

		repo.On("FindByStatus", ctx, "todo", uID).Return([]model.Task{}, nil).Once()
		_, err := svc.GetTasksByStatus(ctx, "todo", uID)
		assert.NoError(t, err)

		repo.On("GetToday", ctx, uID).Return([]model.Task{}, nil).Once()
		_, err = svc.GetTodayTasks(ctx, uID)
		assert.NoError(t, err)
	})
}
