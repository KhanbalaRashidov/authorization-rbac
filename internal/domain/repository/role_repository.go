package repository

import "ms-authz/internal/domain/model"

type RoleRepository interface {
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	GetAllWithPermissions() ([]model.Role, error)
	GetPermissionsByRoleID(id uint) ([]model.Permission, error)
	GetAll() ([]model.Role, error)
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id uint) error
}
