package config

import (
	"os"
	"time"
)

// Config 应用程序统一配置
type Config struct {
	// HTTP 服务配置
	HTTP HTTPConfig `json:"http"`

	// 数据库配置
	Database DatabaseConfig `json:"database"`

	// 日志配置
	Logging LoggingConfig `json:"logging"`

	// OIDC 认证配置
	OIDC OIDCConfig `json:"oidc"`
}

// HTTPConfig HTTP 服务配置
type HTTPConfig struct {
	Addr string `json:"addr"` // 监听地址，如 ":8080"
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN string `json:"dsn"` // SQLite 数据文件路径，如 "data/ddd.db"
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Environment string `json:"environment"` // 运行环境: "development" 或 "production"
	Level       string `json:"level"`       // 日志级别: "debug", "info", "warn", "error"
	OutputPath  string `json:"outputPath"`  // 输出路径: "stdout", "stderr" 或文件路径
}

// OIDCConfig OIDC 认证配置
type OIDCConfig struct {
	IssuerURL     string        `json:"issuerUrl"`     // Keycloak Issuer URL
	ClientID      string        `json:"clientId"`      // Client ID
	ClientSecret  string        `json:"clientSecret"`  // Client Secret
	RedirectURL   string        `json:"redirectUrl"`   // 回调 URL
	Scopes        []string      `json:"scopes"`        // 请求的 scopes
	SessionSecret string        `json:"sessionSecret"` // Session 加密密钥
	SessionTTL    time.Duration `json:"sessionTtl"`    // Session 有效期
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		HTTP: HTTPConfig{
			Addr: ":8080",
		},
		Database: DatabaseConfig{
			DSN: "data/ddd.db",
		},
		Logging: LoggingConfig{
			Environment: "development",
			Level:       "info",
			OutputPath:  "stdout",
		},
		OIDC: OIDCConfig{
			IssuerURL:   "http://localhost:8081/realms/ddd-realm",
			ClientID:    "ddd-app",
			RedirectURL: "http://localhost:8080/oauth2/callback",
			Scopes:      []string{"openid", "profile", "email", "roles"},
			SessionTTL:  24 * time.Hour,
		},
	}
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() Config {
	cfg := DefaultConfig()

	// HTTP 配置
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		cfg.HTTP.Addr = addr
	}

	// 数据库配置
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}

	// 日志配置
	if env := os.Getenv("LOG_ENVIRONMENT"); env != "" {
		cfg.Logging.Environment = env
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Logging.Level = level
	}
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		cfg.Logging.OutputPath = output
	}

	// OIDC 配置
	if issuer := os.Getenv("OIDC_ISSUER_URL"); issuer != "" {
		cfg.OIDC.IssuerURL = issuer
	}
	if clientID := os.Getenv("OIDC_CLIENT_ID"); clientID != "" {
		cfg.OIDC.ClientID = clientID
	}
	if clientSecret := os.Getenv("OIDC_CLIENT_SECRET"); clientSecret != "" {
		cfg.OIDC.ClientSecret = clientSecret
	}
	if redirectURL := os.Getenv("OIDC_REDIRECT_URL"); redirectURL != "" {
		cfg.OIDC.RedirectURL = redirectURL
	}
	if sessionSecret := os.Getenv("OIDC_SESSION_SECRET"); sessionSecret != "" {
		cfg.OIDC.SessionSecret = sessionSecret
	}
	if ttl := os.Getenv("OIDC_SESSION_TTL"); ttl != "" {
		if d, err := time.ParseDuration(ttl); err == nil {
			cfg.OIDC.SessionTTL = d
		}
	}

	return cfg
}
