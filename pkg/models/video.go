package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID       string `gorm:"type:uuid;primaryKey" json:"id"`
	Title    string `gorm:"type:string;not null" json:"title"`
	PublicID string `gorm:"type:varchar(255);not null" json:"public_id"`
	URL      string `gorm:"type:text;not null" json:"url"`
	ModuleID string `gorm:"type:uuid;index" json:"module_id"` // Foreign key referencing Module
}

func (v *Video) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New().String()
	return
}
