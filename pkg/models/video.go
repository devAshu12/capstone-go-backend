package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ID       string `gorm:"type:uuid;primaryKey" json:"video_id"`
	Title    string `gorm:"type:string;not null" json:"title"`
	URL      string `gorm:"type:string;not null" json:"url"`  // URL for the video
	ModuleID string `gorm:"type:uuid;index" json:"module_id"` // Foreign key referencing Module
}

func (v *Video) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New().String()
	return
}
