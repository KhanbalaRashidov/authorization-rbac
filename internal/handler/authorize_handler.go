package handler

import (
	"github.com/gofiber/fiber/v2"
	"ms-authz/internal/infrastructure/mq"
	"ms-authz/internal/service"
	"strings"
)

type AuthorizeHandler struct {
	Auth      *service.AuthService
	RBAC      *service.RBACService
	publisher mq.Publisher
}

func NewAuthorizeHandler(auth *service.AuthService, rbac *service.RBACService, publisher mq.Publisher) *AuthorizeHandler {
	return &AuthorizeHandler{
		Auth:      auth,
		RBAC:      rbac,
		publisher: publisher,
	}
}

func (h *AuthorizeHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/v1/authz/check", h.Authorize)
	app.Post("/api/v1/authz/logout", h.Logout)
	app.Post("/api/v1/authz/logout-all", h.LogoutAll)
}

// Authorize godoc
// @Summary JWT və RBAC yoxlama
// @Description Token JWT ilə doğrulanır. İstəyə əsasən blacklist və RBAC permission da yoxlanır.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param check_jwt query bool false "JWT yoxlanılsın?" default(true)
// @Param check_blacklist query bool false "Blacklist yoxlanılsın?" default(true)
// @Param check_rbac query bool false "RBAC yoxlanılsın?" default(false)
// @Param privilege query string false "RBAC üçün icazə adı (məs: DELETE_USER)"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Privilege is required for RBAC check"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Permission denied"
// @Router /api/v1/authz/check [get]
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

	if claims.UserID != "" && claims.ExpiresAt != nil {
		h.Auth.AddTokenForTracking(token, claims.ExpiresAt.Unix(), claims.UserID, claims.Role)
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

// Logout godoc
// @Summary Logout (Tokeni deaktiv edir)
// @Description İstifadəçi tokenini blackliste əlavə edir (logout əməliyyatı).
// @Tags Authorization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {string} string "Logged out"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/authz/logout [post]
func (h *AuthorizeHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing or invalid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := h.Auth.ParseAndValidate(token, true, false) // Blacklist yoxlamasına ehtiyac yoxdur
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	h.Auth.HandleBlacklistEvent(token, claims.ExpiresAt.Unix())
	h.publishBlacklistEvent(token, claims.ExpiresAt.Unix())

	return c.SendStatus(fiber.StatusOK)
}

type LogoutAllRequest struct {
	UserID string `json:"user_id"`
}

// LogoutAll godoc
// @Summary İstifadəçinin bütün tokenlərini bloklayır
// @Description Verilən `user_id`-yə aid olan bütün JWT-lərin JTI-lərini blackliste əlavə edir və bütün instansiyalara yayır.
// @Tags Authorization
// @Accept json
// @Produce plain
// @Param body body LogoutAllRequest true "Bloklanacaq istifadəçinin ID-si"
// @Success 200 {string} string "All user tokens blacklisted"
// @Failure 400 {string} string "user_id is required"
// @Router /api/v1/authz/logout-all [post]
func (h *AuthorizeHandler) LogoutAll(c *fiber.Ctx) error {
	var req LogoutAllRequest
	if err := c.BodyParser(&req); err != nil || req.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user_id is required")
	}

	// local cache-ə əlavə et
	tokens := h.Auth.GetAllTokensByUser(req.UserID)
	for _, t := range tokens {
		h.Auth.HandleBlacklistEvent(t.Token, t.Exp)
	}

	// RabbitMQ ilə digər instansiyalara yayımlanır
	_ = h.publisher.PublishEvent("auth.tokens.fanout", map[string]any{
		"event":   "TOKEN_BLACKLISTED_ALL",
		"user_id": req.UserID,
	}, []string{
		"blacklist.cache.queue",
		"blacklist.audit.queue",
	})

	return c.SendString("All user tokens blacklisted")
}

func (h *AuthorizeHandler) publishBlacklistEvent(token string, exp int64) {
	event := struct {
		Event string `json:"event"`
		Token string `json:"token"`
		Exp   int64  `json:"exp"`
	}{
		Event: "TOKEN_BLACKLISTED",
		Token: token,
		Exp:   exp,
	}

	_ = h.publisher.PublishEvent("auth.tokens.fanout", event, []string{})
}
