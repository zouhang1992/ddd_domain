package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
)

// NewMeterProvider creates a new OpenTelemetry MeterProvider
func NewMeterProvider(cfg config.Config, logger *zap.Logger, lc fx.Lifecycle) (*sdkmetric.MeterProvider, error) {
	if !cfg.Tracing.Enabled {
		logger.Info("Metrics (OTLP) is disabled")
		return sdkmetric.NewMeterProvider(), nil
	}

	logger.Info("Initializing MeterProvider",
		zap.String("serviceName", cfg.Tracing.ServiceName),
		zap.String("endpoint", cfg.Tracing.Endpoint))

	// Create OTLP metric exporter
	exporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint(cfg.Tracing.Endpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	// Create resource with service name
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.Tracing.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create MeterProvider with periodic reader
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	// Set global MeterProvider
	otel.SetMeterProvider(mp)

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down MeterProvider")
			return mp.Shutdown(ctx)
		},
	})

	logger.Info("MeterProvider initialized successfully")
	return mp, nil
}

// Meter is the global meter instance, initialized after Fx provides NewMeterProvider
var Meter = otel.Meter("ddd-house")
