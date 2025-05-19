package repository

import "ms-authz/internal/domain/model"

type RoleRepository interface {
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	GetAll() ([]model.Role, error)
	Create(role *model.Role) error
	Delete(id uint) error
}
