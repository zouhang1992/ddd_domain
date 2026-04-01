package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
)

// Config 表示日志配置（别名以保持向后兼容）
type Config = config.LoggingConfig

// NewLogger 创建一个新的 zap.Logger 实例
func NewLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
	var zapConfig zap.Config

	switch cfg.Environment {
	case "development":
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case "production":
		zapConfig = zap.NewProductionConfig()
		zapConfig.DisableStacktrace = true
	default:
		zapConfig = zap.NewProductionConfig()
	}

	// 设置日志级别
	switch cfg.Level {
	case "debug":
		zapConfig.Level.SetLevel(zap.DebugLevel)
	case "info":
		zapConfig.Level.SetLevel(zap.InfoLevel)
	case "warn":
		zapConfig.Level.SetLevel(zap.WarnLevel)
	case "error":
		zapConfig.Level.SetLevel(zap.ErrorLevel)
	}

	// 设置输出路径
	if cfg.OutputPath != "stdout" && cfg.OutputPath != "stderr" {
		zapConfig.OutputPaths = []string{cfg.OutputPath}
		zapConfig.ErrorOutputPaths = []string{cfg.OutputPath}
	}

	return zapConfig.Build()
}

// MustNewLogger 是 NewLogger 的便捷函数，会在创建失败时恐慌
func MustNewLogger(cfg config.LoggingConfig) *zap.Logger {
	logger, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	return logger
}

// DefaultConfig 返回默认配置（向后兼容）
func DefaultConfig() Config {
	return Config{
		Environment: "production",
		Level:       "info",
		OutputPath:  "stdout",
	}
}

// Module 提供日志系统的 Uber FX 模块
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg config.Config) config.LoggingConfig {
			return cfg.Logging
		}),
		fx.Provide(NewLogger),
		fx.Invoke(func(l *zap.Logger, cfg config.Config) {
			// 记录应用启动日志
			l.Info("Logging system initialized",
				zap.String("environment", cfg.Logging.Environment),
				zap.String("level", cfg.Logging.Level))
		}),
	)
}
