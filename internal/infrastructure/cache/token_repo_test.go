package cache

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewTokenRepository(t *testing.T) {
	tests := []struct {
		name string
		want *TokenRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenRepo_Add(t *testing.T) {
	type fields struct {
		blacklist  sync.Map
		userTokens sync.Map
	}
	type args struct {
		jti string
		exp int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TokenRepo{
				blacklist:  tt.fields.blacklist,
				userTokens: tt.fields.userTokens,
			}
			r.Add(tt.args.jti, tt.args.exp)
		})
	}
}

func TestTokenRepo_AddWithUser(t *testing.T) {
	type fields struct {
		blacklist  sync.Map
		userTokens sync.Map
	}
	type args struct {
		jti    string
		exp    int64
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TokenRepo{
				blacklist:  tt.fields.blacklist,
				userTokens: tt.fields.userTokens,
			}
			r.AddWithUser(tt.args.jti, tt.args.exp, tt.args.userID)
		})
	}
}

func TestTokenRepo_CleanupExpired(t *testing.T) {
	type fields struct {
		blacklist  sync.Map
		userTokens sync.Map
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TokenRepo{
				blacklist:  tt.fields.blacklist,
				userTokens: tt.fields.userTokens,
			}
			r.CleanupExpired()
		})
	}
}

func TestTokenRepo_GetAllJTIsByUser(t *testing.T) {
	type fields struct {
		blacklist  sync.Map
		userTokens sync.Map
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []TokenInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TokenRepo{
				blacklist:  tt.fields.blacklist,
				userTokens: tt.fields.userTokens,
			}
			if got := r.GetAllJTIsByUser(tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllJTIsByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenRepo_IsBlacklisted(t *testing.T) {
	type fields struct {
		blacklist  sync.Map
		userTokens sync.Map
	}
	type args struct {
		jti string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TokenRepo{
				blacklist:  tt.fields.blacklist,
				userTokens: tt.fields.userTokens,
			}
			if got := r.IsBlacklisted(tt.args.jti); got != tt.want {
				t.Errorf("IsBlacklisted() = %v, want %v", got, tt.want)
			}
		})
	}
}