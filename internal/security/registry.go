package security

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"ticket-api/internal/repository"
)

// SecurityRegistry holds in-memory maps for high performance dynamic PBAC checks
type SecurityRegistry struct {
	mu               sync.RWMutex
	rolesRelations   *repository.RolesRelationsRepository
	rolePermissions  map[int64]map[string]struct{} // role_id -> set of permission_name
	routePermissions map[string][]string           // METHOD:route -> required permissions (ANY match satisfies)
	routeStatus      map[string]bool               // METHOD:route -> active status
}

// NewSecurityRegistry creates and loads the in-memory security matrix
func NewSecurityRegistry(rolesRelations *repository.RolesRelationsRepository) *SecurityRegistry {
	return &SecurityRegistry{
		rolesRelations:   rolesRelations,
		rolePermissions:  make(map[int64]map[string]struct{}),
		routePermissions: make(map[string][]string),
		routeStatus:      make(map[string]bool),
	}
}

// Reload reloads the permissions, roles, and routes mapping from SQLite
func (r *SecurityRegistry) Reload(ctx context.Context) error {
	rolePermRows, err := r.rolesRelations.GetAllRolesWithPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch role permissions: %w", err)
	}

	routePermRows, err := r.rolesRelations.GetAllRoutesWithPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch route permissions: %w", err)
	}

	newRolePerms := make(map[int64]map[string]struct{})
	for _, row := range rolePermRows {
		if _, ok := newRolePerms[row.RoleID]; !ok {
			newRolePerms[row.RoleID] = make(map[string]struct{})
		}
		newRolePerms[row.RoleID][row.PermissionName] = struct{}{}
	}

	newRoutePerms := make(map[string][]string)
	newRouteStatus := make(map[string]bool)

	for _, row := range routePermRows {
		key := fmt.Sprintf("%s:%s", strings.ToUpper(row.Method), row.Route)
		newRouteStatus[key] = row.Status != 0
		if row.PermissionName.Valid && row.PermissionName.String != "" {
			newRoutePerms[key] = append(newRoutePerms[key], row.PermissionName.String)
		}
	}

	r.mu.Lock()
	r.rolePermissions = newRolePerms
	r.routePermissions = newRoutePerms
	r.routeStatus = newRouteStatus
	r.mu.Unlock()

	return nil
}

// CanAccessRoute checks if user roles fulfill route requirements
func (r *SecurityRegistry) CanAccessRoute(roleIDs []int64, method, routePath string) (bool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", strings.ToUpper(method), routePath)

	// Check if route is enabled in registry
	if status, exists := r.routeStatus[key]; exists && !status {
		return false, false // Route disabled
	}

	requiredPerms, hasRules := r.routePermissions[key]
	// If route has no specific permission assigned in DB, allow authenticated access
	if !hasRules || len(requiredPerms) == 0 {
		return true, true
	}

	// Admin shortcut: check if roleIDs contains admin or required permission
	for _, roleID := range roleIDs {
		perms, ok := r.rolePermissions[roleID]
		if !ok {
			continue
		}
		for _, required := range requiredPerms {
			if _, has := perms[required]; has {
				return true, true
			}
		}
	}

	return false, true
}

// GetPermissionsForRoles resolves all permissions for given role IDs
func (r *SecurityRegistry) GetPermissionsForRoles(roleIDs []int64) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	permSet := make(map[string]struct{})
	for _, roleID := range roleIDs {
		if perms, ok := r.rolePermissions[roleID]; ok {
			for p := range perms {
				permSet[p] = struct{}{}
			}
		}
	}

	list := make([]string, 0, len(permSet))
	for p := range permSet {
		list = append(list, p)
	}
	return list
}
