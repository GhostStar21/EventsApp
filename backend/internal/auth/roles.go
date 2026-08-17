package auth

import (
    "context"
    "EventsApp/internal/consts"
)

// GetRoleFromContext returns the role stored in the request context.
func GetRoleFromContext(ctx context.Context) (string, bool) {
    roleVal := ctx.Value("role")
    role, ok := roleVal.(string)
    return role, ok
}

// IsAdmin returns true when the current request role is ADMIN.
func IsAdmin(ctx context.Context) bool {
    role, ok := GetRoleFromContext(ctx)
    return ok && role == string(consts.RoleAdmin)
}
