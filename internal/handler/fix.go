package handler

import (
	"net/http"

	"casdoor-casbin-openbao/internal/database"
	"github.com/labstack/echo/v4"
)

type FixHandler struct{}

func NewFixHandler() *FixHandler {
	return &FixHandler{}
}

// FixCasbinRules fixes ptype field in casbin_rule table
func (h *FixHandler) FixCasbinRules(c echo.Context) error {
	db := database.GetDB()
	if db == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database not connected")
	}

	// Manual SQL fix for ptype
	query1 := "UPDATE casbin_rule SET ptype = 'p' WHERE v2 != '' AND v2 IS NOT NULL"
	query2 := "UPDATE casbin_rule SET ptype = 'g' WHERE v2 = '' OR v2 IS NULL"
	
	result1 := db.Exec(query1)
	result2 := db.Exec(query2)
	
	// Check for errors
	if result1.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Query1 error: "+result1.Error.Error())
	}
	if result2.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Query2 error: "+result2.Error.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ptype field fixed manually",
		"policies_fixed": result1.RowsAffected,
		"roles_fixed": result2.RowsAffected,
	})
}