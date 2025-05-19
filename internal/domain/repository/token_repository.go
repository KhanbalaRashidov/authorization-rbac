package repository

type TokenRepository interface {
	IsBlacklisted(jti string) bool
	Add(jti string, exp int64)
	CleanupExpired()
}
