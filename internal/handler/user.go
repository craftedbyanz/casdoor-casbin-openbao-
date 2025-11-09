package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"casdoor-casbin-openbao/internal/auth"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetProfile returns user profile (protected route)
// GET /api/users/profile
func (h *UserHandler) GetProfile(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          user.GetUserID(),
		"name":        user.Name,
		"display_name": user.DisplayName,
		"email":       user.Email,
		"owner":       user.Owner,
		"roles":       user.Roles,
		"is_admin":    user.IsAdmin,
		"message":     "This is a protected route",
	})
}

// GetUsers returns list of users (admin only)
// GET /api/users
func (h *UserHandler) GetUsers(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "admin access required")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": []map[string]interface{}{
			{
				"id":    user.GetUserID(),
				"name":  user.Name,
				"email": user.Email,
			},
		},
		"message": "This is an admin-only route",
		"admin":   user.Name,
	})
}

// ProtectedResource returns a protected resource
// GET /api/protected
func (h *UserHandler) ProtectedResource(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "This is a protected resource",
		"user":    user.Name,
		"data":    "Sensitive data that requires authentication",
	})
}

// GetSecrets returns secrets from vault (demonstrates cert verification)
// GET /api/secrets
func (h *UserHandler) GetSecrets(c echo.Context) error {
	// Token is verified by AuthMiddleware using Casdoor cert
	// This is where cert verification happens!
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Token verified using Casdoor cert in middleware",
		"user":    user.Name,
		"secrets": []string{"secret1", "secret2", "secret3"},
		"note":    "This endpoint required cert verification via middleware",
	})
}

