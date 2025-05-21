package service

import (
	"errors"
	"ms-authz/internal/domain/repository"
	"ms-authz/pkg/jwtutil"
	"ms-authz/internal/infrastructure/cache"
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

// Token-i yoxlayır və valid claim-ləri qaytarır
func (s *AuthService) ValidateToken(tokenStr string) (*jwtutil.Claims, error) {
	_, kid, err := jwtutil.ParseTokenHeader(tokenStr)
	if err != nil {
		return nil, err
	}

	pubKey, err := s.publicKeyProvider.GetPublicKey(kid)
	if err != nil {
		return nil, err
	}

	parsedClaims, err := jwtutil.VerifyToken(tokenStr, pubKey)
	if err != nil {
		return nil, err
	}

	if s.tokenRepo.IsBlacklisted(parsedClaims.JTI) {
		return nil, errors.New("token is blacklisted")
	}

	return parsedClaims, nil
}

// MQ və ya logout zamanı tokenin `jti`-sini local cache-ə əlavə edir
func (s *AuthService) HandleBlacklistEvent(jti string, exp int64) {
	s.tokenRepo.Add(jti, exp)
}

func (s *AuthService) HandleBlacklistEventWithUser(jti string, exp int64, userID string) {
	s.tokenRepo.AddWithUser(jti, exp, userID)
}


// ParseAndValidate token for jwt + blacklist checks
func (s *AuthService) ParseAndValidate(token string, checkJWT, checkBlacklist bool) (*jwtutil.Claims, error) {
	_, kid, err := jwtutil.ParseTokenHeader(token)
	if err != nil {
		return nil, err
	}

	pubKey, err := s.publicKeyProvider.GetPublicKey(kid)
	if err != nil {
		return nil, err
	}

	parsedClaims, err := jwtutil.VerifyToken(token, pubKey)
	if err != nil {
		if checkJWT {
			return nil, err
		}
	}

	if checkBlacklist && s.tokenRepo.IsBlacklisted(parsedClaims.JTI) {
		return nil, errors.New("token is blacklisted")
	}

	return parsedClaims, nil
}


func (s *AuthService) GetAllJTIsByUser(userID string) []cache.TokenInfo {
	return s.tokenRepo.GetAllJTIsByUser(userID)
}
