package jwtutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type PublicKeyProvider interface {
	GetPublicKey(kid string) (*rsa.PublicKey, error)
}

func ParseTokenHeader(tokenStr string) (*Claims, string, error) {

	//fmt.Println("RAW TOKEN:", tokenStr)
	//fmt.Println("TRIMMED :", strings.TrimSpace(tokenStr))

	parts := strings.Split(tokenStr, ".")
	//fmt.Println("DEBUG: token parts =", len(parts)) //
	if len(parts) != 3 {
		return nil, "", errors.New("invalid token format")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, "", err
	}

	var header struct {
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(decoded, &header); err != nil {
		return nil, "", err
	}

	return &Claims{}, header.Kid, nil
}

func VerifyToken(tokenStr string, pubKey *rsa.PublicKey) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

type fileKeyProvider struct {
	basePath string // Məsələn: /keys/public/
	cache    map[string]*rsa.PublicKey
	mu       sync.RWMutex
}

// PEM fayllar "/keys/public/<kid>.pem" formatındadır
func NewFileKeyProvider(basePath string) PublicKeyProvider {
	return &fileKeyProvider{
		basePath: basePath,
		cache:    make(map[string]*rsa.PublicKey),
	}
}

func (f *fileKeyProvider) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	f.mu.RLock()
	if key, ok := f.cache[kid]; ok {
		f.mu.RUnlock()
		return key, nil
	}
	f.mu.RUnlock()

	filePath := filepath.Join(f.basePath, kid+".pem")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse key: %w", err)
	}

	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	// Cache-ə yaz
	f.mu.Lock()
	f.cache[kid] = rsaKey
	f.mu.Unlock()

	return rsaKey, nil
}

func (f *fileKeyProvider) PreloadKeys() {
	files, _ := os.ReadDir(f.basePath)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pem") {
			kid := strings.TrimSuffix(file.Name(), ".pem")
			_, _ = f.GetPublicKey(kid) // cache-ə yüklə
		}
	}
}
