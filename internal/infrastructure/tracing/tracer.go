package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
)

// Tracer is the global tracer instance
var Tracer trace.Tracer

// NewTracer creates a new OpenTelemetry tracer provider
func NewTracer(cfg config.Config, logger *zap.Logger, lc fx.Lifecycle) (trace.TracerProvider, error) {
	if !cfg.Tracing.Enabled {
		logger.Info("Tracing is disabled")
		return trace.NewNoopTracerProvider(), nil
	}

	logger.Info("Initializing tracer",
		zap.String("serviceName", cfg.Tracing.ServiceName),
		zap.String("endpoint", cfg.Tracing.Endpoint))

	// Create OTLP exporter
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(cfg.Tracing.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
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

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Create global tracer
	Tracer = tp.Tracer(cfg.Tracing.ServiceName)

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down tracer provider")
			return tp.Shutdown(ctx)
		},
	})

	logger.Info("Tracer initialized successfully")
	return tp, nil
}

// GetTracer returns the global tracer
func GetTracer() trace.Tracer {
	if Tracer == nil {
		return otel.Tracer("default")
	}
	return Tracer
}
