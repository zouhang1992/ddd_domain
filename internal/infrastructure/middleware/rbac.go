package middleware

import (
	"net/http"

	"github.com/zouhang1992/ddd_domain/internal/application/auth"
	"go.uber.org/zap"
)

// RBACMiddleware RBAC 中间件
type RBACMiddleware struct {
	rbacService *auth.RBACService
	log         *zap.Logger
}

// NewRBACMiddleware 创建 RBAC 中间件
func NewRBACMiddleware(
	rbacService *auth.RBACService,
	log *zap.Logger,
) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
		log:         log,
	}
}

// RequireRole 要求指定角色的中间件
func (m *RBACMiddleware) RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !m.rbacService.HasRole(claims, role) {
			m.log.Warn("User does not have required role",
				zap.String("user_id", claims.Sub),
				zap.String("required_role", role))
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RequirePermission 要求指定权限的中间件
func (m *RBACMiddleware) RequirePermission(permission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !m.rbacService.HasPermission(claims, permission) {
			m.log.Warn("User does not have required permission",
				zap.String("user_id", claims.Sub),
				zap.String("required_permission", permission))
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RequireAdmin 要求管理员角色的中间件
func (m *RBACMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return m.RequireRole("admin", next)
}
