package dqk

// HasPermission checks whether the user has a specific permission.
// It takes a permission key (e.g., "view:users") and a map of the user's permissions,
// where each key represents a granted permission and the value is typically true.
//
// Returns true if the permission key exists in the map and is set to true.
// Returns false if the permission is missing or explicitly set to false.
//
// Example:
//
//	perms := BuildPermissionSet([]string{"view:users", "edit:posts"})
//	HasPermission("view:users", perms) // true
//	HasPermission("delete:users", perms) // false
func HasPermission(access string, perms map[string]bool) bool {
	val, ok := perms[access]
	return ok && val
}

// IsAccessRestricted determines whether the user lacks both the base and elevated permissions
// required to access a resource.
//
// This is useful in scenarios where you want to allow users with either scoped ("view") or
// broad ("viewAll") permissions to access a resource.
//
// Returns true if the user has neither permission.
// Returns false if the user has at least one of them.
//
// Example:
//
//	IsAccessRestricted("view:article", "viewAll:articles", userPerms)
//	// true → access denied
//	// false → access granted
func IsAccessRestricted(UserPermissions map[string]bool, basePerm, elevatedPerm string) bool {
	return !HasPermission(basePerm, UserPermissions) && !HasPermission(elevatedPerm, UserPermissions)
}

// BuildPermissionSet constructs a map-based set from a slice of permission strings,
// allowing for efficient O(1) permission lookups.
//
// Each string in the input slice is used as a map key with a value of true.
// This is commonly used to preprocess user permission lists into a format suitable
// for fast lookup when checking access.
//
// Example:
//
//	BuildPermissionSet([]string{"view:users", "edit:users"})
//	// → map[string]bool{"view:users": true, "edit:users": true}
func BuildPermissionSet(userPermissions []string) map[string]bool {
	permMap := make(map[string]bool, len(userPermissions))
	for _, perm := range userPermissions {
		permMap[perm] = true
	}
	return permMap
}
