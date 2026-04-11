package tracing

import (
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Module provides the tracing infrastructure as an Uber Fx module
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewTracer),
		fx.Provide(func(tp trace.TracerProvider) trace.Tracer {
			return tp.Tracer("ddd-house")
		}),
	)
}
