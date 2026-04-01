package auth

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
)

// Config OIDC 配置（别名以保持向后兼容）
type Config = config.OIDCConfig

// DefaultConfig 返回默认配置（向后兼容）
func DefaultConfig() Config {
	return Config{
		Scopes:     []string{"openid", "profile", "email", "roles"},
		SessionTTL: 24 * time.Hour,
	}
}

// UserClaims 用户 Claims
type UserClaims struct {
	Sub          string                 `json:"sub"`
	Email        string                 `json:"email"`
	Name         string                 `json:"name"`
	RealmRoles   []string               `json:"realm_roles,omitempty"`
	ResourceRoles map[string][]string    `json:"resource_roles,omitempty"`
	Permissions  []string               `json:"permissions,omitempty"`
	Exp          int64                  `json:"exp"`
	Extra        map[string]any           `json:"-"`
}
