package permissions

import (
	"backend/internal/models"
	"strings"
)

// Role represents a user role in the system.  It's used to control access to different features and resources.
type Role string

const (
	// RoleAdmin Has full access to all system functionalities.
	RoleAdmin Role = "admin"
	// RoleModerator Has elevated privileges compared to a regular user, allowing them to perform moderation tasks.
	RoleModerator Role = "moderator"
	// RoleUser Represents a standard user with basic system access.
	RoleUser Role = "user"
)

// Permission represents a permission for a specific action on a resource
type Permission func(user models.User, data interface{}) bool

// ResourcePermissions represents a map of resource names to their corresponding permissions.
type ResourcePermissions map[string]Permission

// RolesWithPermissions maps roles to their respective permissions for each resource
type RolesWithPermissions map[Role]map[string]ResourcePermissions

// rolesWithPermissions defines a map of roles to their associated permissions. It specifies which actions each role is
// allowed or denied to perform on different resources.
//
// Returns:
//   - RolesWithPermissions: A map containing roles and their corresponding permissions.  Each role maps to a nested map
//     of resources and their permitted actions (e.g., "view," "create," "update," "delete").  Each action maps to a function indicating whether it's allowed or denied (alwaysAllow or alwaysDeny).
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

// HasPermission checks if a user has permission to perform a specific action on a resource.
//
// It iterates through the user's roles, checking if each role grants the required permission.  The permission check
// itself is delegated to a function associated with the specific action.
//
// Parameters:
//   - user (models.User): The user to check permissions for.
//   - resource (string): The resource the action applies to.
//   - action (string): The action being performed.
//   - data (interface{}):  Arbitrary data passed to the permission check function.  This can be used to provide context
//     to the permission check.
//
// Returns:
//   - bool: True if the user has the permission, false otherwise.
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

// isRole checks if the provided role is a valid predefined role.
//
// Parameters:
//   - role (Role): The role to check.
//
// Returns:
//   - bool: True if the role is one of Admin, Moderator, or User; false otherwise.
func isRole(role Role) bool {
	return role == RoleAdmin || role == RoleModerator || role == RoleUser
}
