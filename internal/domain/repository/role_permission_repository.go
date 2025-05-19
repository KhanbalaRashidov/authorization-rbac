package repository

import "ms-authz/internal/domain/model"

type RolePermissionRepository interface {
	GetPermissionsByRoleID(roleID uint) ([]model.Permission, error)
	AddPermission(roleID, permissionID uint) error
	RemovePermission(roleID, permissionID uint) error
	ClearPermissions(roleID uint) error
}
