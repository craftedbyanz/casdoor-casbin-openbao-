package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"casdoor-casbin-openbao/internal/auth"
	"casdoor-casbin-openbao/internal/config"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	config *config.Config
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		config: config.GetConfig(),
	}
}

// LoginRequest represents login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DirectLogin handles direct login with username/password
// POST /api/auth/login
// Body: {"username": "admin", "password": "admin"}
func (h *AuthHandler) DirectLogin(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body: "+err.Error())
	}

	if req.Username == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username and password are required")
	}

	// Login with Casdoor
	token, err := auth.DirectLogin(req.Username, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "login failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"message":      "Login successful. Use the access_token in Authorization header.",
	})
}

// OAuthLogin initiates OAuth login flow
// GET /api/auth/oauth/login
func (h *AuthHandler) OAuthLogin(c echo.Context) error {
	// Generate state for CSRF protection
	state := generateState()

	// Store state in cookie or session (simplified: return in response)
	loginURL := auth.GetLoginURL(state)

	if loginURL == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate login URL: CASDOOR_CLIENT_ID is not configured")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"login_url": loginURL,
		"state":     state,
		"message":   "Redirect to this URL to login",
		"config": map[string]interface{}{
			"endpoint":     h.config.Casdoor.Endpoint,
			"client_id":    h.config.Casdoor.ClientID,
			"organization": h.config.Casdoor.Organization,
			"application":  h.config.Casdoor.Application,
			"redirect_url": h.config.Casdoor.RedirectURL,
		},
	})
}

// Callback handles OAuth callback from Casdoor
// GET /api/auth/callback
func (h *AuthHandler) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	// TODO: Verify state matches original ← Chưa implement!

	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing authorization code")
	}

	// Exchange code for token
	token, err := h.exchangeCodeForToken(code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to exchange token: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token.AccessToken,
		"token_type":   token.TokenType,
		"expires_in":   token.ExpiresIn,
		"state":        state,
		"message":      "Login successful. Use the access_token in Authorization header.",
	})
}

// TokenResponse represents token response from Casdoor
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// exchangeCodeForToken exchanges authorization code for access token
func (h *AuthHandler) exchangeCodeForToken(code string) (*TokenResponse, error) {
	cfg := h.config

	tokenURL := fmt.Sprintf("%s/api/login/oauth/access_token", cfg.Casdoor.Endpoint)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", cfg.Casdoor.ClientID)
	data.Set("client_secret", cfg.Casdoor.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.Casdoor.RedirectURL)

	// Add organization and application if configured
	if cfg.Casdoor.Organization != "" {
		data.Set("organization", cfg.Casdoor.Organization)
	}
	if cfg.Casdoor.Application != "" {
		data.Set("application", cfg.Casdoor.Application)
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// Logout handles logout
// POST /api/auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Logout successful",
	})
}

// GetUserInfo returns current user info from token
// GET /api/auth/me
func (h *AuthHandler) GetUserInfo(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":           user.GetUserID(),
		"name":         user.Name,
		"display_name": user.DisplayName,
		"email":        user.Email,
		"owner":        user.Owner,
		"roles":        user.Roles,
		"is_admin":     user.IsAdmin,
	})
}

// generateState generates a random state string for OAuth
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
