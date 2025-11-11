package casbin

import (
	"fmt"
	"log"

	"casdoor-casbin-openbao/internal/database"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func InitEnforcer() error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, &gormadapter.CasbinRule{}, "casbin_rule")
	if err != nil {
		return fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	Enforcer, err = casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	err = Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("failed to load policy: %w", err)
	}

	// Auto-initialize default policies if DB is empty
	if len(Enforcer.GetPolicy()) == 0 {
		log.Println("No policies found, initializing default policies...")
		if err := initDefaultPoliciesInternal(); err != nil {
			log.Printf("Warning: failed to initialize default policies: %v", err)
		}
	}

	log.Println("Casbin enforcer initialized successfully")
	return nil
}

func GetEnforcer() *casbin.Enforcer {
	return Enforcer
}

// initDefaultPoliciesInternal initializes default policies (internal use)
func initDefaultPoliciesInternal() error {
	if Enforcer == nil {
		return fmt.Errorf("enforcer not initialized")
	}

	// Add default policies
	policies := [][]string{
		// Basic endpoints
		{"admin", "/api/users", "read"},
		{"admin", "/api/admin/*", "write"},
		{"admin", "/api/protected", "read"},
		{"admin", "/api/secrets", "read"},
		{"admin", "/api/auth/me", "read"},
		{"user", "/api/users/profile", "read"},
		{"user", "/api/protected", "read"},
		{"user", "/api/auth/me", "read"},
		
		// Transaction endpoints
		{"admin", "/api/transactions", "read"},           // Admin can see all transactions
		{"admin", "/api/transactions/*", "read"},         // Admin can see specific transactions
		{"user", "/api/transactions/my", "read"},         // User can see own transactions
		{"user", "/api/transactions", "write"},          // User can create transactions
		
		// Order endpoints
		{"admin", "/api/orders", "read"},                 // Admin can see all orders
		{"admin", "/api/orders/*", "read"},               // Admin can see specific orders
		{"admin", "/api/orders/*", "update"},             // Admin can update order status
		{"user", "/api/orders/my", "read"},               // User can see own orders
		{"user", "/api/orders", "write"},                // User can create orders
	}

	for _, policy := range policies {
		_, err := Enforcer.AddPolicy(policy)
		if err != nil {
			return fmt.Errorf("failed to add policy %v: %w", policy, err)
		}
	}

	// Add role assignments
	roles := [][]string{
		{"admin", "admin"},
		{"testuser", "user"},
	}

	for _, role := range roles {
		_, err := Enforcer.AddRoleForUser(role[0], role[1])
		if err != nil {
			return fmt.Errorf("failed to add role %v: %w", role, err)
		}
	}

	log.Println("Default policies initialized successfully")
	return nil
}