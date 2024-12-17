package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Option struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	QuestionID string    `gorm:"type:uuid;not null" json:"question_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	IsCorrect  bool      `gorm:"type:boolean;not null;default:false" json:"is_correct"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;not null" json:"updated_at"`
	Question   Question  `gorm:"foreignKey:QuestionID" json:"question"`
}

func (o *Option) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New().String()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return
}

func (o *Option) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = time.Now()
	return
}
