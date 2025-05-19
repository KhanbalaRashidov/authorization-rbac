package repository

import "ms-authz/internal/domain/model"

type PermissionRepository interface {
	GetByID(id uint) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	GetAll() ([]model.Permission, error)
	Create(permission *model.Permission) error
	Delete(id uint) error
}
