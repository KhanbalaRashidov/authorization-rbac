package service

import (
	"errors"
	"ms-authz/internal/domain/repository"
	"ms-authz/internal/infrastructure/cache"
	"ms-authz/pkg/jwtutil"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	tokenRepo         repository.TokenRepository
	publicKeyProvider jwtutil.PublicKeyProvider
}

func NewAuthService(tokenRepo repository.TokenRepository, keyProvider jwtutil.PublicKeyProvider) *AuthService {
	return &AuthService{
		tokenRepo:         tokenRepo,
		publicKeyProvider: keyProvider,
	}
}

// Yeganə token doğrulama funksiyası – həm JWT, həm Blacklist yoxlayır
func (s *AuthService) Validate(tokenStr string, checkJWT, checkBlacklist bool) (*jwtutil.Claims, error) {
	var claims jwtutil.Claims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid in token header")
		}
		return s.publicKeyProvider.GetPublicKey(kid)
	})
	if err != nil {
		if checkJWT {
			return nil, err
		}
	}

	if !token.Valid && checkJWT {
		return nil, errors.New("invalid token")
	}

	if checkBlacklist && s.tokenRepo.IsBlacklisted(tokenStr) {
		return nil, errors.New("token is blacklisted")
	}

	return &claims, nil
}

// Tokeni Blacklist-ə əlavə et (sadə)
func (s *AuthService) HandleBlacklistEvent(token string, exp int64) {
	s.tokenRepo.Add(token, exp)
}

// Tokeni həm Blacklist-ə, həm user tracking-ə əlavə et
func (s *AuthService) HandleBlacklistEventWithUser(token string, exp int64, userID string, role string) {
	s.tokenRepo.AddWithUser(token, exp, userID, role)
	s.tokenRepo.Add(token, exp)
}

// Sistemdə aktiv olan tokeni izləməyə başla
func (s *AuthService) AddTokenForTracking(token string, exp int64, userID string, role string) {
	s.tokenRepo.AddWithUser(token, exp, userID, role)
}

// İstifadəçiyə aid bütün tokenləri al (admin panel və ya audit üçün)
func (s *AuthService) GetAllTokensByUser(userID string) []cache.TokenInfo {
	return s.tokenRepo.GetAllTokensByUser(userID)
}
