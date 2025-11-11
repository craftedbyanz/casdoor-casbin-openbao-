package main

import (
	"fmt"
	"log"
	"net/http"

	"casdoor-casbin-openbao/internal/auth"
	"casdoor-casbin-openbao/internal/casbin"
	"casdoor-casbin-openbao/internal/config"
	"casdoor-casbin-openbao/internal/database"
	"casdoor-casbin-openbao/internal/handler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional, use environment variables if it doesn't exist
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize config
	config.Init()
	cfg := config.GetConfig()

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Casbin
	if err := casbin.InitEnforcer(); err != nil {
		log.Fatal("Failed to initialize Casbin:", err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Initialize handlers
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	adminHandler := handler.NewAdminHandler()
	debugHandler := handler.NewDebugHandler()
	fixHandler := handler.NewFixHandler()
	microsoftHandler := handler.NewMicrosoftHandler()

	// Serve static files
	e.Static("/", "web")

	// API documentation
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Casdoor Integration Demo API",
			"demo":    "Visit http://localhost:8080 for interactive demo",
			"endpoints": map[string]string{
				"login":     "POST /api/auth/login - Direct login with username/password",
				"oauth":     "GET /api/auth/oauth/login - Get OAuth login URL",
				"microsoft": "GET /api/auth/microsoft/login - Get Microsoft SSO login URL",
				"callback":  "GET /api/auth/callback?code=xxx&state=xxx - OAuth callback",
				"me":        "GET /api/auth/me - Get current user info (requires Bearer token)",
				"profile":   "GET /api/users/profile - Get user profile (requires Bearer token)",
				"protected": "GET /api/protected - Access protected resource (requires Bearer token)",
				"secrets":   "GET /api/secrets - Get secrets (demonstrates cert verification)",
				"users":     "GET /api/users - Get all users (admin only, requires Bearer token)",
			},
		})
	})

	// Auth routes (public)
	authGroup := e.Group("/api/auth")
	// Case 1: Direct login (proxy login)
	authGroup.POST("/login", authHandler.DirectLogin)
	// Case 2: OAuth / OIDC (OpenID Connect) flow (for web apps with frontend)
	authGroup.GET("/oauth/login", authHandler.OAuthLogin)
	// Case 3: Microsoft SSO login
	authGroup.GET("/microsoft/login", microsoftHandler.MicrosoftSSO)
	authGroup.GET("/callback", authHandler.Callback)
	authGroup.POST("/logout", authHandler.Logout)

	// Protected routes (with Casbin authorization)
	protectedGroup := e.Group("/api")
	protectedGroup.Use(auth.AuthMiddleware()) // Authentication
	protectedGroup.Use(casbin.AuthzMiddleware()) // Authorization
	{
		protectedGroup.GET("/auth/me", authHandler.GetUserInfo)
		protectedGroup.GET("/users/profile", userHandler.GetProfile)
		protectedGroup.GET("/protected", userHandler.ProtectedResource)
		protectedGroup.GET("/secrets", userHandler.GetSecrets)
		protectedGroup.GET("/users", userHandler.GetUsers)
	}

	// Admin routes (policy management)
	adminGroup := e.Group("/api/admin")
	adminGroup.Use(auth.AuthMiddleware())
	adminGroup.Use(auth.RequireAdmin()) // Only admin can manage policies
	{
		adminGroup.POST("/init", adminHandler.InitPolicies)
		adminGroup.GET("/policies", adminHandler.GetPolicies)
		adminGroup.POST("/policies", adminHandler.AddPolicy)
		adminGroup.DELETE("/policies", adminHandler.RemovePolicy)
		adminGroup.GET("/roles", adminHandler.GetRoles)
		adminGroup.POST("/roles", adminHandler.AddRole)
		adminGroup.DELETE("/roles", adminHandler.RemoveRole)
		adminGroup.GET("/debug/casbin-rules", debugHandler.GetCasbinRules)
		adminGroup.POST("/debug/fix-casbin", fixHandler.FixCasbinRules)
		adminGroup.POST("/reload-policies", adminHandler.ReloadPolicies)
	}

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", address)
	log.Printf("Casdoor endpoint: %s", cfg.Casdoor.Endpoint)
	log.Printf("Make sure to set CASDOOR_CLIENT_ID and CASDOOR_CLIENT_SECRET environment variables")

	if err := e.Start(address); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}
