package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	QuizID    string    `gorm:"type:uuid;not null" json:"quiz_id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	Points    int       `gorm:"type:int;not null;default:1" json:"points"`
	CreatedAt time.Time `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null" json:"updated_at"`
	Quiz      Quiz      `gorm:"foreignKey:QuizID" json:"quiz"`
	Options   []Option  `gorm:"foreignKey:QuestionID" json:"options"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New().String()
	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()
	return
}

func (q *Question) BeforeUpdate(tx *gorm.DB) (err error) {
	q.UpdatedAt = time.Now()
	return
}
