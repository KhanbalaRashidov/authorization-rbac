package repository

import "ms-authz/internal/infrastructure/cache"

type TokenRepository interface {
	IsBlacklisted(jti string) bool
	Add(jti string, exp int64)
	AddWithUser(token string, exp int64, userID string, role string)
	GetAllJTIsByUser(userID string) []cache.TokenInfo
	GetAllTokensByUser(userID string) []cache.TokenInfo
	CleanupExpired()
}
