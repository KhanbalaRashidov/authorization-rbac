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

// Role endpoints
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

func (h *RBACAdminHandler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.UoW.RoleRepo().GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(roles)
}

func (h *RBACAdminHandler) DeleteRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.UoW.RoleRepo().Delete(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Permission endpoints
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

func (h *RBACAdminHandler) GetPermissions(c *fiber.Ctx) error {
	perms, err := h.UoW.PermissionRepo().GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(perms)
}

func (h *RBACAdminHandler) DeletePermission(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.UoW.PermissionRepo().Delete(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Role-permission assign/remove
func (h *RBACAdminHandler) AssignPermission(c *fiber.Ctx) error {
	roleID, _ := strconv.Atoi(c.Params("roleID"))
	permID, _ := strconv.Atoi(c.Params("permID"))
	if err := h.UoW.RolePermissionRepo().AddPermission(uint(roleID), uint(permID)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *RBACAdminHandler) RemovePermission(c *fiber.Ctx) error {
	roleID, _ := strconv.Atoi(c.Params("roleID"))
	permID, _ := strconv.Atoi(c.Params("permID"))
	if err := h.UoW.RolePermissionRepo().RemovePermission(uint(roleID), uint(permID)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
