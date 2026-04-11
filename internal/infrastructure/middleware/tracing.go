package middleware

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracingHeader is the header name for trace ID
const TracingHeader = "X-Trace-ID"

// Tracing returns an HTTP middleware that creates spans for requests
// and injects trace ID into response headers
func Tracing(tp trace.TracerProvider, serviceName string) func(http.Handler) http.Handler {
	tracer := tp.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from request headers
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Start a new span with HTTP attributes
			spanName := r.Method + " " + r.URL.Path
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.host", r.Host),
					attribute.String("http.scheme", "http"),
					attribute.String("http.proto", r.Proto),
				),
			)
			defer span.End()

			// Get trace ID and add to response header
			spanCtx := span.SpanContext()
			if spanCtx.HasTraceID() {
				w.Header().Set(TracingHeader, spanCtx.TraceID().String())
			}

			// Create a wrapped response writer to capture status code
			wrapped := &statusRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Call the next handler with the new context
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Add status code to span
			span.SetAttributes(
				attribute.Int("http.status_code", wrapped.statusCode),
			)

			// Set span status based on HTTP status code
			if wrapped.statusCode >= 400 {
				span.SetAttributes(attribute.String("error", "true"))
			}
		})
	}
}

// TraceIDFromContext extracts trace ID from context
func TraceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

// SpanIDFromContext extracts span ID from context
func SpanIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}
