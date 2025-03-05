package permissions

import (
	"backend/internal/models"
	"strings"
)

// Role represents the user's role
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleUser      Role = "user"
)

// Permission represents a permission for a specific action on a resource
type Permission func(user models.User, data interface{}) bool

// Permissions map actions to permissions for a specific resource
type ResourcePermissions map[string]Permission

// RolesWithPermissions maps roles to their respective permissions for each resource
type RolesWithPermissions map[Role]map[string]ResourcePermissions

var rolesWithPermissions = RolesWithPermissions{
	RoleAdmin: {
		"novels": {
			"view":   alwaysAllow,
			"create": alwaysAllow,
			"update": alwaysAllow,
		},
		"chapters": {
			"view":   alwaysAllow,
			"create": alwaysAllow,
			"update": alwaysAllow,
			"delete": alwaysAllow,
		},
		"tts": {
			"generate": alwaysAllow,
		},
	},
	RoleModerator: {
		"novels": {
			"view":   alwaysAllow,
			"create": alwaysAllow,
			"update": alwaysAllow,
		},
		"chapters": {
			"view":   alwaysAllow,
			"create": alwaysAllow,
			"update": alwaysAllow,
			"delete": alwaysAllow,
		},
		"tts": {
			"generate": alwaysAllow,
		},
	},
	RoleUser: {
		"novels": {
			"view":   alwaysAllow,
			"create": alwaysDeny,
			"update": alwaysDeny,
		},
		"chapters": {
			"view":   alwaysAllow,
			"create": alwaysDeny,
			"update": alwaysDeny,
			"delete": alwaysDeny,
		},
		"tts": {
			"generate": alwaysAllow,
		},
	},
}

// alwaysAllow is a permission function that always returns true
func alwaysAllow(_ models.User, _ interface{}) bool {
	return true
}

// alwaysDeny is a permission function that always returns false
func alwaysDeny(_ models.User, _ interface{}) bool {
	return false
}

// hasPermission checks if the user has permission for the given action on the resource
func HasPermission(user models.User, resource string, action string, data interface{}) bool {
	roles := strings.Split(user.Roles, ";")
	for _, roleStr := range roles {
		role := Role(roleStr)
		if !isRole(role) {
			continue
		}

		resourcePermissions, ok := rolesWithPermissions[role][resource]
		if !ok {
			continue
		}

		permission, ok := resourcePermissions[action]
		if !ok {
			continue
		}

		if permission(user, data) {
			return true
		}
	}
	return false
}

// isRole checks if a string is a valid Role
func isRole(role Role) bool {
	return role == RoleAdmin || role == RoleModerator || role == RoleUser
}
