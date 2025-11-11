package handler

import (
	"fmt"
	"net/http"
	"strings"

	"casdoor-casbin-openbao/internal/auth"
	"github.com/labstack/echo/v4"
)

type MicrosoftHandler struct{}

func NewMicrosoftHandler() *MicrosoftHandler {
	return &MicrosoftHandler{}
}

// MicrosoftSSO initiates Microsoft SSO login flow
// GET /api/auth/microsoft/login
func (h *MicrosoftHandler) MicrosoftSSO(c echo.Context) error {
	// Generate state for CSRF protection
	state := generateState()

	// Microsoft SSO URL with provider parameter
	microsoftURL := auth.GetMicrosoftLoginURL(state)

	if microsoftURL == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate Microsoft login URL")
	}

	// Check if request wants HTML (from browser)
	accept := c.Request().Header.Get("Accept")
	if strings.Contains(accept, "text/html") {
		// Return HTML with clickable link
		html := fmt.Sprintf(`
			<html>
			<head><title>Microsoft SSO Login</title></head>
			<body>
				<h2>Microsoft SSO Login</h2>
				<p><a href="%s" target="_blank">Click here to login with Microsoft</a></p>
				<p>Or copy this URL:</p>
				<input type="text" value="%s" style="width:100%%;" readonly onclick="this.select()"/>
				<p>State: %s</p>
			</body>
			</html>
		`, microsoftURL, microsoftURL, state)
		return c.HTML(http.StatusOK, html)
	}

	// Return JSON for API calls
	return c.JSON(http.StatusOK, map[string]interface{}{
		"login_url": microsoftURL,
		"state":     state,
		"message":   "Copy login_url and paste in browser",
		"provider":  "microsoft",
	})
}