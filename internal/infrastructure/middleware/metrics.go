package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Instruments holds all OTel metric instruments for the middleware.
// These are created once at startup via NewInstruments and consumed by middleware.
type Instruments struct {
	httpRequestsTotal     metric.Int64Counter
	httpRequestDuration   metric.Float64Histogram
	httpRequestsInFlight  metric.Int64UpDownCounter
	httpRequestSizeBytes  metric.Float64Histogram
	httpResponseSizeBytes metric.Float64Histogram

	commandsExecutedTotal metric.Int64Counter
	queriesExecutedTotal  metric.Int64Counter
	eventsPublishedTotal  metric.Int64Counter

	dbQueriesTotal  metric.Int64Counter
	dbQueryDuration metric.Float64Histogram
}

// NewInstruments creates all OTel metric instruments from the MeterProvider.
func NewInstruments(mp *sdkmetric.MeterProvider) (*Instruments, error) {
	meter := mp.Meter("ddd-house")

	inst := &Instruments{}
	var err error

	// HTTP metrics — names use dots, the OTel Prometheus exporter converts to underscores.
	// Counter "http.requests" → Prometheus: http_requests_total
	// Histogram "http.request.duration" (unit:s) → Prometheus: http_request_duration_seconds

	inst.httpRequestsTotal, err = meter.Int64Counter(
		"http.requests",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	inst.httpRequestDuration, err = meter.Float64Histogram(
		"http.request.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return nil, err
	}

	inst.httpRequestsInFlight, err = meter.Int64UpDownCounter(
		"http.requests.in_flight",
		metric.WithDescription("Number of HTTP requests currently in flight"),
	)
	if err != nil {
		return nil, err
	}

	inst.httpRequestSizeBytes, err = meter.Float64Histogram(
		"http.request.size",
		metric.WithDescription("HTTP request size"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	inst.httpResponseSizeBytes, err = meter.Float64Histogram(
		"http.response.size",
		metric.WithDescription("HTTP response size"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	// Business metrics (reserved for future instrumentation)
	inst.commandsExecutedTotal, err = meter.Int64Counter(
		"commands.executed",
		metric.WithDescription("Total number of commands executed"),
	)
	if err != nil {
		return nil, err
	}

	inst.queriesExecutedTotal, err = meter.Int64Counter(
		"queries.executed",
		metric.WithDescription("Total number of queries executed"),
	)
	if err != nil {
		return nil, err
	}

	inst.eventsPublishedTotal, err = meter.Int64Counter(
		"events.published",
		metric.WithDescription("Total number of domain events published"),
	)
	if err != nil {
		return nil, err
	}

	inst.dbQueriesTotal, err = meter.Int64Counter(
		"db.queries",
		metric.WithDescription("Total number of database queries"),
	)
	if err != nil {
		return nil, err
	}

	inst.dbQueryDuration, err = meter.Float64Histogram(
		"db.query.duration",
		metric.WithDescription("Database query duration"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5),
	)
	if err != nil {
		return nil, err
	}

	return inst, nil
}

// statusRecorder wraps http.ResponseWriter to capture status code and bytes written.
// Used by both Metrics and Tracing middleware (defined once in this package).
type statusRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	n, err := sr.ResponseWriter.Write(b)
	sr.bytesWritten += n
	return n, err
}

// Metrics returns an HTTP middleware that records OTel metrics for each request.
func Metrics(inst *Instruments) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()
			path := r.URL.Path

			inst.httpRequestsInFlight.Add(ctx, 1)
			defer inst.httpRequestsInFlight.Add(ctx, -1)

			// Record request size
			if r.ContentLength > 0 {
				inst.httpRequestSizeBytes.Record(ctx, float64(r.ContentLength),
					metric.WithAttributes(
						attribute.String("method", r.Method),
						attribute.String("path", path),
					),
				)
			}

			sr := &statusRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(sr, r)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(sr.statusCode)

			inst.httpRequestsTotal.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", path),
					attribute.String("status", status),
				),
			)

			inst.httpRequestDuration.Record(ctx, duration,
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", path),
				),
			)

			inst.httpResponseSizeBytes.Record(ctx, float64(sr.bytesWritten),
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", path),
					attribute.String("status", status),
				),
			)
		})
	}
}

// RecordCommand records a command execution metric.
func (inst *Instruments) RecordCommand(ctx context.Context, command string, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	inst.commandsExecutedTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("command", command),
			attribute.String("status", status),
		),
	)
}

// RecordQuery records a query execution metric.
func (inst *Instruments) RecordQuery(ctx context.Context, query string, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	inst.queriesExecutedTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("query", query),
			attribute.String("status", status),
		),
	)
}

// RecordEvent records a domain event publication metric.
func (inst *Instruments) RecordEvent(ctx context.Context, event string) {
	inst.eventsPublishedTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("event", event),
		),
	)
}

// RecordDBQuery records a database query metric.
func (inst *Instruments) RecordDBQuery(ctx context.Context, operation string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	inst.dbQueriesTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("operation", operation),
			attribute.String("status", status),
		),
	)
	inst.dbQueryDuration.Record(ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("operation", operation),
		),
	)
}
