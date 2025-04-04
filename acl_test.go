package dynamicquerykit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name        string
		access      string
		permissions map[string]bool
		expected    bool
	}{
		{
			name:   "has permission",
			access: "view:users",
			permissions: map[string]bool{
				"view:users": true,
			},
			expected: true,
		},
		{
			name:   "no permission",
			access: "viewAll:users",
			permissions: map[string]bool{
				"view:users": true,
			},
			expected: false,
		},
		{
			name:   "no permission explicit",
			access: "viewAll:users",
			permissions: map[string]bool{
				"view:users":    true,
				"viewAll:users": false,
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasPermission(tt.access, tt.permissions)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestIsAccessRestricted(t *testing.T) {
	tests := []struct {
		name           string
		baseAccess     string
		elevatedAccess string
		permissions    map[string]bool
		expected       bool
	}{
		{
			name:           "has base permission",
			baseAccess:     "view:users",
			elevatedAccess: "viewAll:users",
			permissions: map[string]bool{
				"view:users": true,
			},
			expected: false,
		},
		{
			name:           "has elevated permission",
			baseAccess:     "view:users",
			elevatedAccess: "viewAll:users",
			permissions: map[string]bool{
				"viewAll:users": true,
			},
			expected: false,
		},
		{
			name:           "has both permission",
			baseAccess:     "view:users",
			elevatedAccess: "viewAll:users",
			permissions: map[string]bool{
				"viewAll:users": true,
				"view:users":    true,
			},
			expected: false,
		},
		{
			name:           "no permission",
			baseAccess:     "view:users",
			elevatedAccess: "viewAll:users",
			permissions: map[string]bool{
				"view:cars":    true,
				"viewAll:cars": true,
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAccessRestricted(tt.permissions, tt.baseAccess, tt.elevatedAccess)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestBuildPermissionSet(t *testing.T) {
	tests := []struct {
		name            string
		permissionSlice []string
		expected        map[string]bool
	}{
		{
			name: "slice to map",
			permissionSlice: []string{
				"view:users",
				"viewAll:users",
				"view:cars",
				"viewAll:parts",
			},
			expected: map[string]bool{
				"view:users":    true,
				"viewAll:users": true,
				"view:cars":     true,
				"viewAll:parts": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPermissionSet(tt.permissionSlice)
			assert.Equal(t, tt.expected, got)
		})
	}
}
