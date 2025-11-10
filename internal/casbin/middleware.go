package casbin

import (
	"net/http"

	"casdoor-casbin-openbao/internal/auth"
	"github.com/labstack/echo/v4"
)

// AuthzMiddleware checks authorization using Casbin
func AuthzMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := auth.GetUserFromContext(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
			}

			// Get request info
			path := c.Request().URL.Path
			method := c.Request().Method

			// Convert HTTP method to action
			action := getActionFromMethod(method)
			
			// Check permission: enforce(subject, object, action)
			allowed, err := GetEnforcer().Enforce(user.Name, path, action)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "authorization check failed")
			}

			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			return next(c)
		}
	}
}

func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "write"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "read"
	}
}