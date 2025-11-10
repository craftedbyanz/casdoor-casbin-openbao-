package casbin

import "fmt"

// InitDefaultPolicies initializes default policies and roles (API endpoint)
func InitDefaultPolicies() error {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("enforcer not initialized")
	}

	// Clear existing policies
	enforcer.ClearPolicy()

	// Use internal function to add policies
	return initDefaultPoliciesInternal()
}