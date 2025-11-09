package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"casdoor-casbin-openbao/internal/config"
)

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response from Casdoor
type LoginResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   string `json:"data"` // JWT token
	Data2  string `json:"data2"`
	Data3  bool   `json:"data3"`
}

// DirectLogin logs in user directly with Casdoor API
// This is simpler than OAuth flow and suitable for API-only backends
// Username can be either "username" or "owner/username" format
func DirectLogin(username, password string) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", fmt.Errorf("config not initialized")
	}

	loginURL := fmt.Sprintf("%s/api/login", cfg.Casdoor.Endpoint)

	// Casdoor login API accepts JSON format
	// Username can be either "username" or "owner/username" format
	// If username contains "/", use as is (already has owner/name format)
	// Otherwise, use username as is (Casdoor will use organization from request)
	loginReq := map[string]interface{}{
		"application": cfg.Casdoor.Application,
		"username":    username, // Use username as is, Casdoor will handle organization
		"password":    password,
		"type":        "token",
	}

	// Add organization if configured (Casdoor needs this to find user)
	if cfg.Casdoor.Organization != "" {
		loginReq["organization"] = cfg.Casdoor.Organization
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request: %w", err)
	}

	fmt.Println("loginReq: ", string(jsonData))

	resp, err := http.Post(loginURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to request login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	if loginResp.Status != "ok" {
		return "", fmt.Errorf("login failed: %s", loginResp.Msg)
	}

	if loginResp.Data == "" {
		return "", fmt.Errorf("login failed: no token returned")
	}

	return loginResp.Data, nil
}
