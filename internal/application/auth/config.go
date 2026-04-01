package auth

import "time"

// Config OIDC 配置
type Config struct {
	IssuerURL     string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	Scopes        []string
	SessionSecret string
	SessionTTL    time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Scopes:     []string{"openid", "profile", "email", "roles"},
		SessionTTL: 24 * time.Hour,
	}
}
