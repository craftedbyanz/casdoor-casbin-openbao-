package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"casdoor-casbin-openbao/internal/config"
)

// AuthMiddleware validates JWT tokens from Casdoor
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			token := parts[1]
			claims, err := VerifyToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
			}

			// Store user info in context
			c.Set("user", claims)
			c.Set("user_id", claims.GetUserID())
			c.Set("user_name", claims.Name)
			c.Set("user_email", claims.Email)
			c.Set("is_admin", claims.IsAdmin)

			return next(c)
		}
	}
}

// GetUserFromContext retrieves user claims from echo context
func GetUserFromContext(c echo.Context) (*CasdoorClaims, bool) {
	user, ok := c.Get("user").(*CasdoorClaims)
	return user, ok
}

// RequireAdmin middleware that requires admin role
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := GetUserFromContext(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
			}

			if !user.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "admin access required")
			}

			return next(c)
		}
	}
}

// GetLoginURL generates the OAuth login URL
func GetLoginURL(state string) string {
	cfg := config.GetConfig()
	if cfg == nil {
		return ""
	}

	// Check if client_id is configured
	if cfg.Casdoor.ClientID == "" {
		return ""
	}

	// Casdoor OAuth endpoint format: /login/oauth/authorize
	// Include organization and application if configured
	params := url.Values{}
	params.Set("client_id", cfg.Casdoor.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", cfg.Casdoor.RedirectURL)
	params.Set("scope", "read")
	params.Set("state", state)
	
	// Add organization and application if configured (some Casdoor versions require these)
	if cfg.Casdoor.Organization != "" {
		params.Set("organization", cfg.Casdoor.Organization)
	}
	if cfg.Casdoor.Application != "" {
		params.Set("application", cfg.Casdoor.Application)
	}
	
	loginURL := fmt.Sprintf("%s/login/oauth/authorize?%s",
		cfg.Casdoor.Endpoint,
		params.Encode(),
	)
	
	return loginURL
}

