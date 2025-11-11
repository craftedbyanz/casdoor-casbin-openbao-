package auth

import (
	"fmt"
	"net/url"

	"casdoor-casbin-openbao/internal/config"
)

// GetMicrosoftLoginURL generates Microsoft SSO login URL via Casdoor
func GetMicrosoftLoginURL(state string) string {
	cfg := config.GetConfig()
	if cfg == nil {
		return ""
	}

	if cfg.Casdoor.ClientID == "" {
		return ""
	}

	// Microsoft SSO via Casdoor with provider parameter
	params := url.Values{}
	params.Set("client_id", cfg.Casdoor.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", cfg.Casdoor.RedirectURL) // This is for Casdoor → Backend
	params.Set("scope", "read")
	params.Set("state", state)
	params.Set("provider", "microsoft-provider") // Microsoft provider name in Casdoor

	if cfg.Casdoor.Organization != "" {
		params.Set("organization", cfg.Casdoor.Organization)
	}
	if cfg.Casdoor.Application != "" {
		params.Set("application", cfg.Casdoor.Application)
	}

	microsoftURL := fmt.Sprintf("%s/login/oauth/authorize?%s",
		cfg.Casdoor.Endpoint,
		params.Encode(),
	)

	return microsoftURL
}