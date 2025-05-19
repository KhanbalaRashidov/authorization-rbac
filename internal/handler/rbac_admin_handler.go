package handler

import (
	"github.com/gofiber/fiber/v2"
	"ms-authz/internal/domain/model"
	"ms-authz/internal/domain/repository"
	"strconv"
)

type RBACAdminHandler struct {
	UoW repository.UnitOfWork
}

func NewRBACAdminHandler(uow repository.UnitOfWork) *RBACAdminHandler {
	return &RBACAdminHandler{UoW: uow}
}

func (h *RBACAdminHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/roles", h.CreateRole)
	app.Get("/roles", h.GetRoles)
	app.Delete("/roles/:id", h.DeleteRole)

	app.Post("/permissions", h.CreatePermission)
	app.Get("/permissions", h.GetPermissions)
	app.Delete("/permissions/:id", h.DeletePermission)

	app.Post("/roles/:roleID/permissions/:permID", h.AssignPermission)
	app.Delete("/roles/:roleID/permissions/:permID", h.RemovePermission)
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
	return c.JSON(roles)
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
	return c.JSON(perms)
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
	return c.SendStatus(fiber.StatusNoContent)
}
