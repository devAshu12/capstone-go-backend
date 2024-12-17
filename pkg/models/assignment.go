package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Assignment struct {
	ID          string `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string `gorm:"type:string;not null" json:"title"`
	Description string `gorm:"type:varchar(255);not null" json:"description"`

	ModuleID string `gorm:"type:uuid;index" json:"module_id"`
	Module   Module `gorm:"foreignKey:ModuleID" json:"module"`

	CourseID string `gorm:"type:uuid;index" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID" json:"course"`

	Deadline    time.Time    `gorm:"type:timestamp;not null" json:"deadline"`
	FacultyID   string       `gorm:"type:uuid;index;foreignKey:ID;references:users" json:"faculty_id"`
	Submissions []Submission `gorm:"foreignKey:AssignmentID" json:"submissions"`
}

func (a *Assignment) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()
	return
}
