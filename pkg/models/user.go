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
	ID         string   `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName  string   `gorm:"type:string;not null" json:"first_name"`
	SecondName string   `gorm:"type:string;not null" json:"second_name"`
	Email      string   `gorm:"type:string;not null;unique" json:"-"`
	Password   string   `gorm:"type:string;not null" json:"-"`
	Role       RoleType `gorm:"type:string;not null" json:"-"`
	Valid      bool     `gorm:"type:bool;default:true" json:"-"`
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
