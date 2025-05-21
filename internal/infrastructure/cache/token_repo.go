package cache

import (
	"sync"
	"time"
)

type TokenInfo struct {
	Token string // JWT özü
	Exp   int64  // Expiration timestamp
}

type TokenRepo struct {
	blacklist  sync.Map // map[jti]exp_timestamp
	userTokens sync.Map
}

func NewTokenRepository() *TokenRepo {
	return &TokenRepo{}
}

func (r *TokenRepo) IsBlacklisted(token string) bool {
	value, ok := r.blacklist.Load(token)
	if !ok {
		return false
	}
	exp := value.(int64)
	return time.Now().Unix() < exp
}

func (r *TokenRepo) Add(token string, exp int64) {
	r.blacklist.Store(token, exp)
}

func (r *TokenRepo) AddWithUser(token string, exp int64, userID string, role string) {
	val, _ := r.userTokens.LoadOrStore(userID, []TokenInfo{})
	tokenList := val.([]TokenInfo)

	for _, t := range tokenList {
		if t.Token == token {
			return // artıq var
		}
	}

	tokenList = append(tokenList, TokenInfo{
		Token: token,
		Exp:   exp,
	})
	r.userTokens.Store(userID, tokenList)
}

func (r *TokenRepo) GetAllJTIsByUser(userID string) []TokenInfo {
	val, ok := r.userTokens.Load(userID)
	if !ok {
		return nil
	}
	return val.([]TokenInfo)
}

func (r *TokenRepo) GetAllTokensByUser(userID string) []TokenInfo {
	val, ok := r.userTokens.Load(userID)
	if !ok {
		return nil
	}
	return val.([]TokenInfo)
}

func (r *TokenRepo) CleanupExpired() {
	now := time.Now().Unix()

	// Blacklist təmizlənməsi
	r.blacklist.Range(func(key, value any) bool {
		exp := value.(int64)
		if exp < now {
			r.blacklist.Delete(key)
		}
		return true
	})

	// userTokens təmizlənməsi
	r.userTokens.Range(func(userKey, tokenList any) bool {
		userID := userKey.(string)
		tokens := tokenList.([]TokenInfo)

		var validTokens []TokenInfo
		for _, t := range tokens {
			if t.Exp > now {
				validTokens = append(validTokens, t)
			}
		}

		if len(validTokens) == 0 {
			r.userTokens.Delete(userID)
		} else {
			r.userTokens.Store(userID, validTokens)
		}
		return true
	})
}
