package handler

import (
	"net/http"

	"casdoor-casbin-openbao/internal/database"
	"github.com/labstack/echo/v4"
)

type DebugHandler struct{}

func NewDebugHandler() *DebugHandler {
	return &DebugHandler{}
}

type CasbinRule struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	PType string `json:"ptype"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
}

// GetCasbinRules returns raw casbin_rule table data
func (h *DebugHandler) GetCasbinRules(c echo.Context) error {
	db := database.GetDB()
	if db == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database not connected")
	}

	var rules []CasbinRule
	if err := db.Table("casbin_rule").Find(&rules).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to query casbin_rule: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total": len(rules),
		"rules": rules,
	})
}