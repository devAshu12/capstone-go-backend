package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quiz struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string     `gorm:"type:text;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;not null" json:"updated_at"`
	DueDate     time.Time  `gorm:"type:timestamp" json:"due_date"`
	CourseID    string     `gorm:"type:uuid;not null" json:"course_id"`
	ModuleID    *string    `gorm:"type:uuid" json:"module_id"`
	Active      bool       `gorm:"type:boolean;not null;default:false" json:"active"`
	IsFinal     bool       `gorm:"type:boolean;not null;default:false" json:"is_final"`
	Course      Course     `gorm:"foreignKey:CourseID" json:"course"`
	Module      *Module    `gorm:"foreignKey:ModuleID" json:"module"`
	Questions   []Question `gorm:"foreignKey:QuizID" json:"questions"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New().String()
	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()
	return
}

func (q *Quiz) BeforeUpdate(tx *gorm.DB) (err error) {
	q.UpdatedAt = time.Now()
	return
}

// if question length is 0, set active to false
func (q *Quiz) BeforeSave(tx *gorm.DB) (err error) {
	if len(q.Questions) == 0 {
		q.Active = false
	}
	return
}

// if question length is 0, set active to false
func (q *Quiz) AfterSave(tx *gorm.DB) (err error) {
	if len(q.Questions) == 0 {
		q.Active = false
	}
	return
}
