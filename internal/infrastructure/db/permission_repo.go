package db

import (
	"gorm.io/gorm"
	"ms-authz/internal/domain/model"
)

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) GetByID(id uint) (*model.Permission, error) {
	var p model.Permission
	err := r.db.First(&p, id).Error
	return &p, err
}

func (r *PermissionRepo) GetByName(name string) (*model.Permission, error) {
	var p model.Permission
	err := r.db.Where("name = ?", name).First(&p).Error
	return &p, err
}

func (r *PermissionRepo) GetAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepo) Create(p *model.Permission) error {
	return r.db.Create(p).Error
}

func (r *PermissionRepo) Update(role *model.Permission) error {
	return r.db.Save(role).Error
}

func (r *PermissionRepo) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

func (r *PermissionRepo) GetAllWithRoles() ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.Preload("Roles").Find(&perms).Error
	return perms, err
}

func (r *PermissionRepo) GetRolesByPermissionID(id uint) ([]model.Role, error) {
	var perm model.Permission
	err := r.db.Preload("Roles").First(&perm, id).Error
	return perm.Roles, err
}
