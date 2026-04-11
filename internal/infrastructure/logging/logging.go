package logging

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
)

// globalLogger 是全局的 logger 实例
var globalLogger *zap.Logger

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

	logger, err := zapConfig.Build()
	if err == nil {
		// 设置全局 logger
		globalLogger = logger
	}
	return logger, err
}

// SetGlobalLogger 设置全局 logger
func SetGlobalLogger(logger *zap.Logger) {
	globalLogger = logger
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

// WithTrace 从 context 中提取 trace ID 和 span ID，并添加到 logger 中
func WithTrace(ctx context.Context, logger *zap.Logger) *zap.Logger {
	// 如果没有传入 logger，使用全局 logger
	if logger == nil {
		logger = globalLogger
	}
	// 如果全局 logger 也是 nil，创建一个默认的
	if logger == nil {
		logger = zap.L()
	}

	traceID := middleware.TraceIDFromContext(ctx)
	spanID := middleware.SpanIDFromContext(ctx)

	if traceID != "" {
		logger = logger.With(zap.String("trace_id", traceID))
	}
	if spanID != "" {
		logger = logger.With(zap.String("span_id", spanID))
	}

	return logger
}

// Ctx 从 context 中提取 trace ID 和 span ID，返回带有这些字段的全局 logger
func Ctx(ctx context.Context) *zap.Logger {
	return WithTrace(ctx, nil)
}

// Module 提供日志系统的 Uber FX 模块
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg config.Config) config.LoggingConfig {
			return cfg.Logging
		}),
		fx.Provide(NewLogger),
		fx.Invoke(func(l *zap.Logger, cfg config.Config) {
			// 设置全局 logger
			SetGlobalLogger(l)
			// 记录应用启动日志
			l.Info("Logging system initialized",
				zap.String("environment", cfg.Logging.Environment),
				zap.String("level", cfg.Logging.Level))
		}),
	)
}
