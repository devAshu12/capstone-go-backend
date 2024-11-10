package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleType string

const (
	SuperAdminDev RoleType = "super_admin_dev"
	SuperAdmin    RoleType = "super_admin"
	Faculty       RoleType = "faculty"
	Student       RoleType = "student"
)

type User struct {
	gorm.Model
	ID       string   `gorm:"type:uuid; primaryKey" json:"user_id"`
	Email    string   `gorm:"type:string; not null" json:"email"`
	Password string   `gorm:"type:string; not null" json:"password"`
	Role     RoleType `gorm:"type:string; not null" json:"role"`
	Valid    bool     `gorm:"type:bool; default:true" json:"valid"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	switch u.Role {
	case SuperAdminDev, SuperAdmin, Faculty, Student:
		return nil
	default:
		return fmt.Errorf("invalid role: %s", u.Role)
	}
}
