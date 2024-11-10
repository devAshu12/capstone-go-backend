package models

import "gorm.io/gorm"

type ModuleVideoProgress struct {
	gorm.Model
	UserID     string  `gorm:"type:uuid;index" json:"user_id"`      // Foreign key referencing User
	ModuleID   string  `gorm:"type:uuid;index" json:"module_id"`    // Foreign key referencing Video
	Progress   float64 `gorm:"type:float;not null" json:"progress"` // Progress percentage (0-100)
	Completion bool    `gorm:"type:bool; default:false" json:"completion"`
}
