package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID    string  `gorm:"type:uuid; primaryKey" json:"id"`
	Title string  `gorm:"type:string; not null" json:"title"`
	Price float64 `gorm:"type:decimal(10,2); not null" json:"price"`

	TotalModules     int `gorm:"default:0" json:"total_modules"`
	TotalVideos      int `gorm:"default:0" json:"total_videos"`
	TotalHours       int `gorm:"default:0" json:"total_hours"`
	TotalQuizs       int `gorm:"default:0" json:"total_quizs"`
	TotalAssignments int `gorm:"default:0" json:"total_assignments"`

	// Foreign key to User
	UserRefer string `gorm:"type:uuid;index" json:"user_refer"` // Field to store User ID

	// Faculty relationship with the User model
	Faculty User `gorm:"foreignKey:UserRefer;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"faculty"`

	//one course has many modules
	Modules     []Module     `gorm:"foreignKey:CourseID" json:"modules"`         // Relationship to Modules
	Assignments []Assignment `gorm:"foreignKey:CourseID" json:"assignments"`     // Relationship to Assignments
	Quizzes     []Quiz       `gorm:"foreignKey:CourseID" json:"quizzes"`         // Relationship to Quizzes
	Videos      []Video      `gorm:"foreignKey:CourseID" json:"videos"`          // Relationship to Videos
	Students    []User       `gorm:"many2many:course_students;" json:"students"` // Relationship to Students
}

func (c *Course) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
