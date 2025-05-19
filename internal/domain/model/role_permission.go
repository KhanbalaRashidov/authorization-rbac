package model

import "gorm.io/gorm"

type RolePermission struct {
	gorm.Model
	RoleID       uint `gorm:"not null;index"`
	PermissionID uint `gorm:"not null;index"`
}
