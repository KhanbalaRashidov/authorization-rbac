package repository

type UnitOfWork interface {
	RoleRepo() RoleRepository
	PermissionRepo() PermissionRepository
	RolePermissionRepo() RolePermissionRepository
	UserRepo() UserRepository
}
