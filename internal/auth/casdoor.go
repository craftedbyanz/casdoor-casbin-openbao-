package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"casdoor-casbin-openbao/internal/config"
)

var (
	certCache     *rsa.PublicKey
	certCacheLock sync.RWMutex
)

// CasdoorClaims represents the JWT claims from Casdoor
type CasdoorClaims struct {
	Owner       string   `json:"owner"`
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email"`
	ID          string   `json:"id"`
	Roles       []string `json:"roles"`
	IsAdmin     bool     `json:"isAdmin"`
	jwt.RegisteredClaims
}

// GetUserID returns the user ID from claims (prefer Subject, fallback to ID)
func (c *CasdoorClaims) GetUserID() string {
	if c.Subject != "" {
		return c.Subject
	}
	return c.ID
}

// GetPublicKey fetches the public key from Casdoor
func GetPublicKey() (*rsa.PublicKey, error) {
	certCacheLock.RLock()
	if certCache != nil {
		certCacheLock.RUnlock()
		return certCache, nil
	}
	certCacheLock.RUnlock()

	certCacheLock.Lock()
	defer certCacheLock.Unlock()

	// Double check
	if certCache != nil {
		return certCache, nil
	}

	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}

	// Fetch JWKS from Casdoor (public endpoint, no auth required)
	jwksURL := fmt.Sprintf("%s/.well-known/jwks", cfg.Casdoor.Endpoint)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS: %w", err)
	}

	// Parse JWKS JSON
	var jwks struct {
		Keys []struct {
			Kty string   `json:"kty"`
			Kid string   `json:"kid"`
			Use string   `json:"use"`
			X5c []string `json:"x5c"` // Certificate chain (base64)
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	if len(jwks.Keys) == 0 {
		return nil, errors.New("no keys found in JWKS")
	}

	// Get the first key (usually there's only one)
	key := jwks.Keys[0]
	if len(key.X5c) == 0 {
		return nil, errors.New("no certificate found in JWKS key")
	}

	// Decode base64 certificate (x5c[0] is the certificate)
	certDER, err := base64.StdEncoding.DecodeString(key.X5c[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode certificate: %w", err)
	}

	// Parse certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("certificate is not RSA")
	}

	certCache = pubKey
	return pubKey, nil
}

// VerifyToken verifies a JWT token from Casdoor
func VerifyToken(tokenString string) (*CasdoorClaims, error) {
	publicKey, err := GetPublicKey()
	if err != nil {
		fmt.Println("publicKey-err: ", err)
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CasdoorClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		fmt.Println("token-err: ", err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CasdoorClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

