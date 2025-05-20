package handler

import (
	"github.com/gofiber/fiber/v2"
	"ms-authz/internal/service"
	"strings"
	"github.com/rabbitmq/amqp091-go"
	"encoding/json"
	"log"
)

type AuthorizeHandler struct {
	Auth *service.AuthService
	RBAC *service.RBACService
	RabbitMQCh  *amqp091.Channel
}

func NewAuthorizeHandler(auth *service.AuthService, rbac *service.RBACService, rabbitCh *amqp091.Channel) *AuthorizeHandler {
	return &AuthorizeHandler{
		Auth:        auth,
		RBAC:        rbac,
		RabbitMQCh:  rabbitCh,
	}
}


func (h *AuthorizeHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/authorize", h.Authorize)
	app.Post("/logout", h.Logout)
	app.Post("/admin/block-token", h.BlockToken)
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

// Logout godoc
// @Summary Logout (Tokeni deaktiv edir)
// @Description İstifadəçi tokenini blackliste əlavə edir (logout əməliyyatı).
// @Tags Authorization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {string} string "Logged out"
// @Failure 401 {string} string "Unauthorized"
// @Router /logout [post]
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

	h.Auth.HandleBlacklistEvent(claims.JTI, claims.ExpiresAt.Unix())
	h.publishBlacklistEvent(claims.JTI, claims.ExpiresAt.Unix())

	return c.SendStatus(fiber.StatusOK)
}


type BlockTokenRequest struct {
	JTI string `json:"jti"`
	Exp int64  `json:"exp"` // Unix timestamp
}

// BlockToken godoc
// @Summary Admin token bloklama
// @Description Admin tərəfindən manual olaraq JWT `jti` və `exp`-ə əsasən tokenin blackliste əlavə olunması
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body BlockTokenRequest true "JTI və Exp göndər"
// @Success 200 {string} string "Token blocked"
// @Failure 400 {string} string "Validation error"
// @Router /admin/block-token [post]
func (h *AuthorizeHandler) BlockToken(c *fiber.Ctx) error {
	var req BlockTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}
	if req.JTI == "" || req.Exp == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Missing JTI or Exp")
	}

	h.Auth.HandleBlacklistEvent(req.JTI, req.Exp)
	h.publishBlacklistEvent(req.JTI, req.Exp)
	return c.SendString("Token blocked")
}


func (h *AuthorizeHandler) publishBlacklistEvent(jti string, exp int64) {
	event := struct {
		Event string `json:"event"`
		JTI   string `json:"jti"`
		Exp   int64  `json:"exp"`
	}{
		Event: "TOKEN_BLACKLISTED",
		JTI:   jti,
		Exp:   exp,
	}

	body, _ := json.Marshal(event)
	err := h.RabbitMQCh.Publish(
		"auth.tokens.fanout", // Exchange
		"",                   // Routing key (ignored for fanout)
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("❌ Failed to publish TOKEN_BLACKLISTED event:", err)
	} else {
		log.Println("📤 TOKEN_BLACKLISTED event published")
	}
}
