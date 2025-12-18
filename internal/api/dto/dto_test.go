package dto

import (
	"testing"
	"time"
	"todo-list/internal/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskDTOs(t *testing.T) {
	t.Run("TaskRequestDTO_Check", func(t *testing.T) {
		req := TaskRequestDTO{
			Title: "Test Title",
		}
		assert.Equal(t, "Test Title", req.Title)
	})

	t.Run("ToTaskResponseDTO_Mapping", func(t *testing.T) {
		id := uuid.New()
		now := time.Now()
		task := model.Task{
			ID:        id,
			Title:     "Model Title",
			Content:   "Content",
			Status:    "todo",
			Priority:  "high",
			Archived:  false,
			CreatedAt: now,
			UpdatedAt: now,
			Tags: []model.Tag{
				{Name: "Go"},
				{Name: "Backend"},
			},
		}

		response := ToTaskResponseDTO(task)

		assert.Equal(t, id.String(), response.ID)
		assert.Equal(t, "Model Title", response.Title)
		assert.ElementsMatch(t, []string{"Go", "Backend"}, response.Tags)
		assert.Equal(t, now, response.CreatedAt)
	})
}
