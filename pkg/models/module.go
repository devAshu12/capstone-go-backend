package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Module struct {
	ID               string `gorm:"type:uuid;primaryKey" json:"id"`
	Title            string `gorm:"type:string;not null" json:"title"`
	TotalVideos      int    `gorm:"default:0" json:"total_videos"`
	TotalHours       int    `gorm:"default:0" json:"total_hours"`
	TotalQuizs       int    `gorm:"default:0" json:"total_quizs"`
	TotalAssignments int    `gorm:"default:0" json:"total_assignments"`

	CourseID string  `gorm:"type:uuid;index" json:"course_id"`  // Foreign key referencing Course
	Videos   []Video `gorm:"foreignKey:ModuleID" json:"videos"` // Relationship to Videos
}

func (m *Module) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New().String()
	return
}
