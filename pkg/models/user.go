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
	ID         string   `gorm:"type:uuid;primaryKey" json:"user_id"`
	FirstName  string   `gorm:"type:string;not null" json:"first_name"`
	SecondName string   `gorm:"type:string;not null" json:"second_name"`
	Email      string   `gorm:"type:string;not null;unique" json:"email"`
	Password   string   `gorm:"type:string;not null" json:"password"`
	Role       RoleType `gorm:"type:string;not null" json:"role"`
	Valid      bool     `gorm:"type:bool;default:true" json:"valid"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if !isValidRole(u.Role) {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return
}

// isValidRole checks if a given role is valid
func isValidRole(role RoleType) bool {
	switch role {
	case SuperAdminDev, SuperAdmin, Faculty, Student:
		return true
	}
	return false
}
