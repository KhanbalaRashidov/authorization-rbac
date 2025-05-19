package db

import (
	"gorm.io/gorm"
	"ms-authz/internal/domain/model"
)

type RolePermissionRepo struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepo {
	return &RolePermissionRepo{db: db}
}

func (r *RolePermissionRepo) GetPermissionsByRoleID(roleID uint) ([]model.Permission, error) {
	var role model.Role
	err := r.db.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (r *RolePermissionRepo) AddPermission(roleID, permissionID uint) error {
	role := model.Role{Model: gorm.Model{ID: roleID}}
	perm := model.Permission{Model: gorm.Model{ID: permissionID}}
	return r.db.Model(&role).Association("Permissions").Append(&perm)
}

func (r *RolePermissionRepo) RemovePermission(roleID, permissionID uint) error {
	role := model.Role{Model: gorm.Model{ID: roleID}}
	perm := model.Permission{Model: gorm.Model{ID: permissionID}}
	return r.db.Model(&role).Association("Permissions").Delete(&perm)
}

func (r *RolePermissionRepo) ClearPermissions(roleID uint) error {
	role := model.Role{Model: gorm.Model{ID: roleID}}
	return r.db.Model(&role).Association("Permissions").Clear()
}
