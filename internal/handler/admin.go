package handler

import (
	"net/http"

	"casdoor-casbin-openbao/internal/casbin"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

type PolicyRequest struct {
	Subject string `json:"subject"`
	Object  string `json:"object"`
	Action  string `json:"action"`
}

type RoleRequest struct {
	User string `json:"user"`
	Role string `json:"role"`
}

// GetPolicies returns all policies
func (h *AdminHandler) GetPolicies(c echo.Context) error {
	policies := casbin.GetPolicies()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"policies": policies,
	})
}

// AddPolicy adds a new policy
func (h *AdminHandler) AddPolicy(c echo.Context) error {
	var req PolicyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := casbin.AddPolicy(req.Subject, req.Object, req.Action); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "policy added successfully",
	})
}

// RemovePolicy removes a policy
func (h *AdminHandler) RemovePolicy(c echo.Context) error {
	var req PolicyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := casbin.RemovePolicy(req.Subject, req.Object, req.Action); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "policy removed successfully",
	})
}

// GetRoles returns all role assignments
func (h *AdminHandler) GetRoles(c echo.Context) error {
	roles := casbin.GetRoles()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": roles,
	})
}

// AddRole assigns a role to user
func (h *AdminHandler) AddRole(c echo.Context) error {
	var req RoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := casbin.AddRoleForUser(req.User, req.Role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "role assigned successfully",
	})
}

// RemoveRole removes a role from user
func (h *AdminHandler) RemoveRole(c echo.Context) error {
	var req RoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := casbin.DeleteRoleForUser(req.User, req.Role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "role removed successfully",
	})
}

// InitPolicies initializes default policies and roles
func (h *AdminHandler) InitPolicies(c echo.Context) error {
	if err := casbin.InitDefaultPolicies(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "default policies initialized successfully",
	})
}

// ReloadPolicies reloads policies from database
func (h *AdminHandler) ReloadPolicies(c echo.Context) error {
	if err := casbin.ReloadPolicies(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "policies reloaded successfully",
	})
}