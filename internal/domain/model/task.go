package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Status    string         `gorm:"type:varchar(50);default:'todo'" json:"status"`
	Priority  string         `gorm:"type:varchar(50);default:'medium'" json:"priority"`
	Tags      []Tag          `gorm:"many2many:task_tags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tags"`
	DueDate   *time.Time     `json:"due_date"`
	Archived  bool           `gorm:"default:false" json:"archived"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;type:varchar(100)" json:"name"`
}

// BeforeCreate GORM hook to set UUID
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
