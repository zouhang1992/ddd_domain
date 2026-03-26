package logging

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 表示日志配置
type Config struct {
	// Environment 表示运行环境，支持 "development" 或 "production"
	Environment string `json:"environment"`
	// Level 表示日志级别，支持 "debug", "info", "warn", "error"
	Level string `json:"level"`
	// OutputPath 表示日志输出路径，默认为标准输出
	OutputPath string `json:"outputPath"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Environment: "production",
		Level:       "info",
		OutputPath:  "stdout",
	}
}

// NewLogger 创建一个新的 zap.Logger 实例
func NewLogger(cfg Config) (*zap.Logger, error) {
	var config zap.Config

	switch cfg.Environment {
	case "development":
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case "production":
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
	default:
		config = zap.NewProductionConfig()
	}

	// 设置日志级别
	switch cfg.Level {
	case "debug":
		config.Level.SetLevel(zap.DebugLevel)
	case "info":
		config.Level.SetLevel(zap.InfoLevel)
	case "warn":
		config.Level.SetLevel(zap.WarnLevel)
	case "error":
		config.Level.SetLevel(zap.ErrorLevel)
	}

	// 设置输出路径
	if cfg.OutputPath != "stdout" && cfg.OutputPath != "stderr" {
		config.OutputPaths = []string{cfg.OutputPath}
		config.ErrorOutputPaths = []string{cfg.OutputPath}
	}

	return config.Build()
}

// MustNewLogger 是 NewLogger 的便捷函数，会在创建失败时恐慌
func MustNewLogger(cfg Config) *zap.Logger {
	logger, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	return logger
}

// Module 提供日志系统的 Uber FX 模块
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() Config {
			// 从环境变量或默认配置获取
			env := os.Getenv("LOG_ENVIRONMENT")
			if env == "" {
				env = "production"
			}

			level := os.Getenv("LOG_LEVEL")
			if level == "" {
				level = "info"
			}

			return Config{
				Environment: env,
				Level:       level,
				OutputPath:  "stdout",
			}
		}),
		fx.Provide(NewLogger),
		fx.Invoke(func(l *zap.Logger) {
			// 记录应用启动日志
			l.Info("Logging system initialized",
				zap.String("environment", os.Getenv("LOG_ENVIRONMENT")),
				zap.String("level", os.Getenv("LOG_LEVEL")))
		}),
	)
}
