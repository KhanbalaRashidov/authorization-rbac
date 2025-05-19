package handler

import (
	"github.com/gofiber/fiber/v2"
	"ms-authz/internal/service"
	"strings"
)

type AuthorizeHandler struct {
	Auth *service.AuthService
	RBAC *service.RBACService
}

func NewAuthorizeHandler(auth *service.AuthService, rbac *service.RBACService) *AuthorizeHandler {
	return &AuthorizeHandler{Auth: auth, RBAC: rbac}
}

func (h *AuthorizeHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/authorize", h.Authorize)
}

// Authorize godoc
// @Summary JWT və RBAC yoxlama
// @Description Token JWT ilə doğrulanır. İstəyə əsasən blacklist və RBAC permission da yoxlanır.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param check_jwt query bool true "JWT imzası yoxlansın? (default: true)"
// @Param check_blacklist query bool true "Token blacklistedir? (default: true)"
// @Param check_rbac query bool false "RBAC permission yoxlansın? (default: false)"
// @Param privilege query string false "RBAC üçün icazə adı (məs: DELETE_USER)"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Privilege is required for RBAC check"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Permission denied"
// @Router /authorize [get]
func (h *AuthorizeHandler) Authorize(c *fiber.Ctx) error {
	// 1. Token oxu
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing or invalid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. Query parametrləri oxu
	checkJWT := c.QueryBool("check_jwt", true)
	checkBlacklist := c.QueryBool("check_blacklist", true)
	checkRBAC := c.QueryBool("check_rbac", false)
	privilege := c.Query("privilege", "")

	// 3. Token parse və yoxlama
	claims, err := h.Auth.ParseAndValidate(token, checkJWT, checkBlacklist)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// 4. RBAC yoxlama
	if checkRBAC {
		if privilege == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Privilege is required for RBAC check")
		}
		if !h.RBAC.HasPermission(claims.Role, privilege) {
			return fiber.NewError(fiber.StatusForbidden, "Permission denied")
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
