package db

import (
	"errors"
	"gorm.io/gorm"
	"ms-authz/internal/domain/repository"
)

type GormUnitOfWork struct {
	db                 *gorm.DB
	tx                 *gorm.DB
	roleRepo           repository.RoleRepository
	permissionRepo     repository.PermissionRepository
	rolePermissionRepo repository.RolePermissionRepository
	userRepo           repository.UserRepository
}

func NewUnitOfWork(db *gorm.DB) repository.UnitOfWork {
	return &GormUnitOfWork{db: db}
}

func (u *GormUnitOfWork) begin() *gorm.DB {
	if u.tx != nil {
		return u.tx
	}
	u.tx = u.db.Begin()
	return u.tx
}

// Transaction Start (optional public method)
func (u *GormUnitOfWork) Begin() {
	u.tx = u.db.Begin()
}

// Transaction Commit
func (u *GormUnitOfWork) Commit() error {
	if u.tx == nil {
		return errors.New("no active transaction")
	}
	err := u.tx.Commit().Error
	u.tx = nil
	return err
}

// Transaction Rollback
func (u *GormUnitOfWork) Rollback() error {
	if u.tx == nil {
		return errors.New("no active transaction")
	}
	err := u.tx.Rollback().Error
	u.tx = nil
	return err
}

// RoleRepo getter
func (u *GormUnitOfWork) RoleRepo() repository.RoleRepository {
	if u.roleRepo == nil {
		u.roleRepo = NewRoleRepository(u.getDB())
	}
	return u.roleRepo
}

// PermissionRepo getter
func (u *GormUnitOfWork) PermissionRepo() repository.PermissionRepository {
	if u.permissionRepo == nil {
		u.permissionRepo = NewPermissionRepository(u.getDB())
	}
	return u.permissionRepo
}

// RolePermissionRepo getter
func (u *GormUnitOfWork) RolePermissionRepo() repository.RolePermissionRepository {
	if u.rolePermissionRepo == nil {
		u.rolePermissionRepo = NewRolePermissionRepository(u.getDB())
	}
	return u.rolePermissionRepo
}

// UserRepo getter
func (u *GormUnitOfWork) UserRepo() repository.UserRepository {
	if u.userRepo == nil {
		u.userRepo = NewUserRepository(u.getDB())
	}
	return u.userRepo
}

// Internal helper for choosing correct DB (with or without transaction)
func (u *GormUnitOfWork) getDB() *gorm.DB {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}
