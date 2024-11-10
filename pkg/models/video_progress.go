package models

import "gorm.io/gorm"

type UserVideoProgress struct {
	gorm.Model
	UserID   string  `gorm:"type:uuid;index" json:"user_id"`      // Foreign key referencing User
	VideoID  string  `gorm:"type:uuid;index" json:"video_id"`     // Foreign key referencing Video
	Progress float64 `gorm:"type:float;not null" json:"progress"` // Progress percentage (0-100)

	TimeSpent  float32 `gorm:"type:float;not null" json:"time_spent"`
	Completion bool    `gorm:"type:bool; default:false" json:"completion"`
}
