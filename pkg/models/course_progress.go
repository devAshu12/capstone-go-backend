package models

import "gorm.io/gorm"

type CourseVideoProgress struct {
	gorm.Model
	UserID string `gorm:"type:uuid;index" json:"user_id"` // Foreign key referencing User
	User   User   `gorm:"foreignKey:UserID" json:"user"`

	CourseID string `gorm:"type:uuid;index" json:"course_id"` // Foreign key referencing Video
	Course   Course `gorm:"foreignKey:CourseID" json:"course"`

	Progress float64 `gorm:"type:float;not null" json:"progress"` // Progress percentage (0-100)

	Completion bool `gorm:"type:bool; default:false" json:"completion"`
}
