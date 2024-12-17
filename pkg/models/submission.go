package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Submission struct {
	ID           string     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       string     `gorm:"type:uuid;not null" json:"user_id"`
	AssignmentID string     `gorm:"type:uuid;not null" json:"assignment_id"`
	SubmittedAt  time.Time  `gorm:"type:timestamp;not null" json:"submitted_at"`
	FileURL      *string    `gorm:"type:text" json:"file_url"`
	User         User       `gorm:"foreignKey:UserID" json:"user"`
	Assignment   Assignment `gorm:"foreignKey:AssignmentID" json:"assignment"`
}

func (s *Submission) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New().String()
	s.SubmittedAt = time.Now()
	return
}
