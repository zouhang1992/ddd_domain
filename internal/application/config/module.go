package config

import "go.uber.org/fx"

// Module 提供配置模块
var Module = fx.Options(
	fx.Provide(LoadFromEnv),
)
