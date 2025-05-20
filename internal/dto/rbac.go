package dto

type RoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PermissionDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RoleWithPermissionsDTO struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Permissions []PermissionDTO `json:"permissions"`
}

type PermissionWithRolesDTO struct {
	ID    uint       `json:"id"`
	Name  string     `json:"name"`
	Roles []RoleDTO  `json:"roles"`
}
