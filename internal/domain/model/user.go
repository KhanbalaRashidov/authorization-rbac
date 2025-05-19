package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:100;not null"`
	Email    string `gorm:"uniqueIndex;size:150"`
	RoleID   uint
	Role     Role
}
