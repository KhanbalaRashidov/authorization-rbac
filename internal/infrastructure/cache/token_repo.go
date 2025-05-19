package cache

import (
	"sync"
	"time"
)

type TokenRepo struct {
	blacklist sync.Map // map[jti]exp_timestamp
}

func NewTokenRepository() *TokenRepo {
	return &TokenRepo{}
}

func (r *TokenRepo) IsBlacklisted(jti string) bool {
	value, ok := r.blacklist.Load(jti)
	if !ok {
		return false
	}
	exp := value.(int64)
	return time.Now().Unix() < exp
}

func (r *TokenRepo) Add(jti string, exp int64) {
	r.blacklist.Store(jti, exp)
}

func (r *TokenRepo) CleanupExpired() {
	now := time.Now().Unix()
	r.blacklist.Range(func(key, value any) bool {
		exp := value.(int64)
		if exp < now {
			r.blacklist.Delete(key)
		}
		return true
	})
}
