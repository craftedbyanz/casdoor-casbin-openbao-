package casbin

import "fmt"

// AddPolicy adds a policy rule
func AddPolicy(sub, obj, act string) error {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("enforcer not initialized")
	}

	added, err := enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		return fmt.Errorf("failed to add policy: %w", err)
	}

	if !added {
		return fmt.Errorf("policy already exists")
	}

	return nil
}

// RemovePolicy removes a policy rule
func RemovePolicy(sub, obj, act string) error {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("enforcer not initialized")
	}

	removed, err := enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}

	if !removed {
		return fmt.Errorf("policy not found")
	}

	return nil
}

// AddRoleForUser assigns a role to user
func AddRoleForUser(user, role string) error {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("enforcer not initialized")
	}

	added, err := enforcer.AddRoleForUser(user, role)
	if err != nil {
		return fmt.Errorf("failed to add role: %w", err)
	}

	if !added {
		return fmt.Errorf("role assignment already exists")
	}

	return nil
}

// DeleteRoleForUser removes a role from user
func DeleteRoleForUser(user, role string) error {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("enforcer not initialized")
	}

	deleted, err := enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if !deleted {
		return fmt.Errorf("role assignment not found")
	}

	return nil
}

// GetPolicies returns all policies
func GetPolicies() [][]string {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return nil
	}
	return enforcer.GetPolicy()
}

// GetRoles returns all role assignments
func GetRoles() [][]string {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return nil
	}
	return enforcer.GetGroupingPolicy()
}

// ReloadPolicies reloads policies from database
func ReloadPolicies() error {
	enforcer := GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("enforcer not initialized")
	}

	return enforcer.LoadPolicy()
}