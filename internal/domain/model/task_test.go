package model

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskModel(t *testing.T) {
	id := uuid.New()
	task := Task{
		ID:    id,
		Title: "Model Test",
	}

	assert.Equal(t, id, task.ID)
	assert.Equal(t, "Model Test", task.Title)
}
