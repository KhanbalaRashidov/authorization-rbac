package repository

import "ms-authz/internal/domain/model"

type UserRepository interface {
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetAll() ([]model.User, error)
	Create(user *model.User) error
	Delete(id uint) error
}
