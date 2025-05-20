package handler

import (
	"github.com/gofiber/fiber/v2"
	"ms-authz/internal/domain/model"
	"ms-authz/internal/domain/repository"
	"ms-authz/internal/service"
	"ms-authz/internal/dto"
	"strconv"
)

type RBACAdminHandler struct {
	UoW repository.UnitOfWork
	RBAC *service.RBACService
}

func NewRBACAdminHandler(uow repository.UnitOfWork, rbacService *service.RBACService) *RBACAdminHandler {
	return &RBACAdminHandler{UoW: uow, RBAC: rbacService}
}

func (h *RBACAdminHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/roles", h.CreateRole)
	app.Get("/roles", h.GetRoles)
	app.Delete("/roles/:id", h.DeleteRole)
	app.Post("/roles/:roleID/permissions/:permID", h.AssignPermission)
	app.Delete("/roles/:roleID/permissions/:permID", h.RemovePermission)
	app.Get("/roles-with-permissions", h.GetRolesWithPermissions)
	app.Get("/roles/:id/permissions", h.GetPermissionsByRoleID)
	app.Put("/roles/:id", h.UpdateRole)

	app.Post("/permissions", h.CreatePermission)
	app.Get("/permissions", h.GetPermissions)
	app.Delete("/permissions/:id", h.DeletePermission)
	app.Get("/permissions-with-roles", h.GetPermissionsWithRoles)
	app.Get("/permissions/:id/roles", h.GetRolesByPermissionID)
	app.Put("/permissions/:id", h.UpdatePermission)

}

// CreateRole godoc
// @Summary Yeni rol yaradır
// @Tags Role
// @Accept json
// @Produce json
// @Param role body model.Role true "Yeni rol"
// @Success 200 {object} model.Role
// @Failure 400 {string} string "Invalid body"
// @Failure 500 {string} string "Server error"
// @Router /roles [post]
func (h *RBACAdminHandler) CreateRole(c *fiber.Ctx) error {
	var role model.Role
	if err := c.BodyParser(&role); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}
	if err := h.UoW.RoleRepo().Create(&role); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_ROLE_CREATED", map[string]any{
		"role_id": role.ID,
		"role_name": role.Name,
	})

	return c.JSON(role)
}

// GetRoles godoc
// @Summary Mövcud bütün rolları qaytarır
// @Tags Role
// @Produce json
// @Success 200 {array} model.Role
// @Failure 500 {string} string "Server error"
// @Router /roles [get]
func (h *RBACAdminHandler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.UoW.RoleRepo().GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var result []dto.RoleDTO
	for _, r := range roles {
		result = append(result, dto.RoleDTO{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return c.JSON(result)
}


// UpdateRole godoc
// @Summary Mövcud rolu yeniləyir
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body model.Role true "Yenilənmiş rol məlumatı"
// @Success 200 {object} model.Role
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Role not found"
// @Failure 500 {string} string "Server error"
// @Router /roles/{id} [put]
func (h *RBACAdminHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	var updated model.Role
	if err := c.BodyParser(&updated); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}

	role, err := h.UoW.RoleRepo().GetByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Role not found")
	}

	role.Name = updated.Name

	if err := h.UoW.RoleRepo().Update(role); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_ROLE_UPDATED", map[string]any{
		"role_id": role.ID,
		"new_name": updated.Name,
	})

	return c.JSON(role)
}

// DeleteRole godoc
// @Summary Rolu ID-yə görə silir
// @Tags Role
// @Param id path int true "Role ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Server error"
// @Router /roles/{id} [delete]
func (h *RBACAdminHandler) DeleteRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.UoW.RoleRepo().Delete(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_ROLE_DELETED", map[string]any{
		"role_id": id,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

// CreatePermission godoc
// @Summary Yeni permission yaradır
// @Tags Permission
// @Accept json
// @Produce json
// @Param permission body model.Permission true "Yeni permission"
// @Success 200 {object} model.Permission
// @Failure 400 {string} string "Invalid body"
// @Failure 500 {string} string "Server error"
// @Router /permissions [post]
func (h *RBACAdminHandler) CreatePermission(c *fiber.Ctx) error {
	var p model.Permission
	if err := c.BodyParser(&p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}
	if err := h.UoW.PermissionRepo().Create(&p); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_PERMISSION_CREATED", map[string]any{
		"perm_id": p.ID,
		"perm_name": p.Name,
	})

	return c.JSON(p)
}

// GetPermissions godoc
// @Summary Bütün permission-ları qaytarır
// @Tags Permission
// @Produce json
// @Success 200 {array} model.Permission
// @Failure 500 {string} string "Server error"
// @Router /permissions [get]
func (h *RBACAdminHandler) GetPermissions(c *fiber.Ctx) error {
	perms, err := h.UoW.PermissionRepo().GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var result []dto.PermissionDTO
	for _, p := range perms {
		result = append(result, dto.PermissionDTO{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return c.JSON(result)
}


// UpdatePermission godoc
// @Summary Mövcud permission-u yeniləyir
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param permission body model.Permission true "Yenilənmiş permission məlumatı"
// @Success 200 {object} model.Permission
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Permission not found"
// @Failure 500 {string} string "Server error"
// @Router /permissions/{id} [put]
func (h *RBACAdminHandler) UpdatePermission(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	var updated model.Permission
	if err := c.BodyParser(&updated); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}

	perm, err := h.UoW.PermissionRepo().GetByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Permission not found")
	}

	perm.Name = updated.Name

	if err := h.UoW.PermissionRepo().Update(perm); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_PERMISSION_UPDATED", map[string]any{
		"perm_id": perm.ID,
		"new_name": updated.Name,
	})

	return c.JSON(perm)
}

// DeletePermission godoc
// @Summary Permission-u ID ilə silir
// @Tags Permission
// @Param id path int true "Permission ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Server error"
// @Router /permissions/{id} [delete]
func (h *RBACAdminHandler) DeletePermission(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.UoW.PermissionRepo().Delete(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_PERMISSION_DELETED", map[string]any{
		"perm_id": id,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

// AssignPermission godoc
// @Summary Role-a permission təyin edir
// @Tags Role-Permission
// @Param roleID path int true "Role ID"
// @Param permID path int true "Permission ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Server error"
// @Router /roles/{roleID}/permissions/{permID} [post]
func (h *RBACAdminHandler) AssignPermission(c *fiber.Ctx) error {
	roleID, _ := strconv.Atoi(c.Params("roleID"))
	permID, _ := strconv.Atoi(c.Params("permID"))
	if err := h.UoW.RolePermissionRepo().AddPermission(uint(roleID), uint(permID)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_PERMISSION_ASSIGNED", map[string]any{
		"role_id": roleID,
		"perm_id": permID,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

// RemovePermission godoc
// @Summary Role-dan permission silir
// @Tags Role-Permission
// @Param roleID path int true "Role ID"
// @Param permID path int true "Permission ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Server error"
// @Router /roles/{roleID}/permissions/{permID} [delete]
func (h *RBACAdminHandler) RemovePermission(c *fiber.Ctx) error {
	roleID, _ := strconv.Atoi(c.Params("roleID"))
	permID, _ := strconv.Atoi(c.Params("permID"))
	if err := h.UoW.RolePermissionRepo().RemovePermission(uint(roleID), uint(permID)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.RBAC.PublishCacheEvent("RBAC_PERMISSION_REMOVED", map[string]any{
		"role_id": roleID,
		"perm_id": permID,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

// GetRolesWithPermissions godoc
// @Summary Rolları və onlara bağlı permission-ları qaytarır
// @Tags Role
// @Produce json
// @Success 200 {array} model.Role
// @Failure 500 {string} string "Server error"
// @Router /roles-with-permissions [get]
func (h *RBACAdminHandler) GetRolesWithPermissions(c *fiber.Ctx) error {
	roles, err := h.UoW.RoleRepo().GetAllWithPermissions()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var result []dto.RoleWithPermissionsDTO
	for _, r := range roles {
		roleDTO := dto.RoleWithPermissionsDTO{
			ID:   r.ID,
			Name: r.Name,
		}
		for _, p := range r.Permissions {
			roleDTO.Permissions = append(roleDTO.Permissions, dto.PermissionDTO{
				ID:   p.ID,
				Name: p.Name,
			})
		}
		result = append(result, roleDTO)
	}

	return c.JSON(result)
}


// GetPermissionsWithRoles godoc
// @Summary Permission-ları və aid olduqları rolları qaytarır
// @Tags Permission
// @Produce json
// @Success 200 {array} model.Permission
// @Failure 500 {string} string "Server error"
// @Router /permissions-with-roles [get]
func (h *RBACAdminHandler) GetPermissionsWithRoles(c *fiber.Ctx) error {
	perms, err := h.UoW.PermissionRepo().GetAllWithRoles()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var result []dto.PermissionWithRolesDTO
	for _, p := range perms {
		permDTO := dto.PermissionWithRolesDTO{
			ID:   p.ID,
			Name: p.Name,
		}
		for _, r := range p.Roles {
			permDTO.Roles = append(permDTO.Roles, dto.RoleDTO{
				ID:   r.ID,
				Name: r.Name,
			})
		}
		result = append(result, permDTO)
	}

	return c.JSON(result)
}

// GetPermissionsByRoleID godoc
// @Summary Verilmiş role ID üçün permission-ları qaytarır
// @Tags Role
// @Param id path int true "Role ID"
// @Produce json
// @Success 200 {array} model.Permission
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Server error"
// @Router /roles/{id}/permissions [get]
func (h *RBACAdminHandler) GetPermissionsByRoleID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	perms, err := h.UoW.RoleRepo().GetPermissionsByRoleID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var result []dto.PermissionDTO
	for _, p := range perms {
		result = append(result, dto.PermissionDTO{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return c.JSON(result)
}

// GetRolesByPermissionID godoc
// @Summary Verilmiş permission ID üçün aid olduğu rolları qaytarır
// @Tags Permission
// @Param id path int true "Permission ID"
// @Produce json
// @Success 200 {array} model.Role
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Server error"
// @Router /permissions/{id}/roles [get]
func (h *RBACAdminHandler) GetRolesByPermissionID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	roles, err := h.UoW.PermissionRepo().GetRolesByPermissionID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var result []dto.RoleDTO
	for _, r := range roles {
		result = append(result, dto.RoleDTO{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return c.JSON(result)
}
