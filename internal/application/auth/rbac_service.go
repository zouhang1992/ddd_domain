package auth

import "go.uber.org/zap"

// RBACService RBAC 服务
type RBACService struct {
	log *zap.Logger
}

// NewRBACService 创建 RBAC 服务
func NewRBACService(log *zap.Logger) *RBACService {
	return &RBACService{log: log}
}

// HasRole 检查用户是否拥有指定角色
func (s *RBACService) HasRole(claims *UserClaims, role string) bool {
	if claims == nil {
		return false
	}

	// 检查 realm roles
	for _, r := range claims.RealmRoles {
		if r == role {
			return true
		}
	}

	// 检查 resource roles
	for _, roles := range claims.ResourceRoles {
		for _, r := range roles {
			if r == role {
				return true
			}
		}
	}

	return false
}

// HasPermission 检查用户是否拥有指定权限
func (s *RBACService) HasPermission(claims *UserClaims, permission string) bool {
	if claims == nil {
		return false
	}

	// 管理员拥有所有权限
	if s.HasRole(claims, "admin") {
		return true
	}

	// 检查权限列表
	for _, p := range claims.Permissions {
		if p == permission {
			return true
		}
	}

	return false
}

// GetRoles 获取用户所有角色
func (s *RBACService) GetRoles(claims *UserClaims) []string {
	if claims == nil {
		return nil
	}

	var roles []string
	roles = append(roles, claims.RealmRoles...)

	for _, rs := range claims.ResourceRoles {
		roles = append(roles, rs...)
	}

	// 去重
	seen := make(map[string]bool)
	result := make([]string, 0, len(roles))
	for _, r := range roles {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}

	return result
}

// GetPermissions 获取用户所有权限
func (s *RBACService) GetPermissions(claims *UserClaims) []string {
	if claims == nil {
		return nil
	}
	return claims.Permissions
}

// IsAdmin 检查用户是否是管理员
func (s *RBACService) IsAdmin(claims *UserClaims) bool {
	return s.HasRole(claims, "admin")
}
