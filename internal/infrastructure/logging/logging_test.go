package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "stdout", cfg.OutputPath)
}

func TestNewLogger(t *testing.T) {
	t.Run("development environment", func(t *testing.T) {
		cfg := Config{
			Environment: "development",
			Level:       "debug",
			OutputPath:  "stdout",
		}
		logger, err := NewLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		logger.Sync()
	})

	t.Run("production environment", func(t *testing.T) {
		cfg := Config{
			Environment: "production",
			Level:       "info",
			OutputPath:  "stdout",
		}
		logger, err := NewLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		logger.Sync()
	})

	t.Run("warn level", func(t *testing.T) {
		cfg := Config{
			Environment: "production",
			Level:       "warn",
			OutputPath:  "stdout",
		}
		logger, err := NewLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		logger.Sync()
	})

	t.Run("error level", func(t *testing.T) {
		cfg := Config{
			Environment: "production",
			Level:       "error",
			OutputPath:  "stdout",
		}
		logger, err := NewLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		logger.Sync()
	})
}

func TestMustNewLogger(t *testing.T) {
	cfg := DefaultConfig()
	logger := MustNewLogger(cfg)
	assert.NotNil(t, logger)
	logger.Sync()
}

func TestLoggingOutput(t *testing.T) {
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	logger.Debug("debug message", zap.String("key", "value"))
	logger.Info("info message", zap.String("key", "value"))
	logger.Warn("warn message", zap.String("key", "value"))
	logger.Error("error message", zap.String("key", "value"))

	assert.Equal(t, 4, observedLogs.Len())
	assert.Equal(t, "debug message", observedLogs.All()[0].Message)
	assert.Equal(t, "info message", observedLogs.All()[1].Message)
	assert.Equal(t, "warn message", observedLogs.All()[2].Message)
	assert.Equal(t, "error message", observedLogs.All()[3].Message)
}
