package db

import (
	"gorm.io/gorm"
	"ms-authz/internal/domain/model"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *RoleRepo) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *RoleRepo) GetAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepo) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}
