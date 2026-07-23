package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP 指标
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently in flight",
		},
	)

	httpRequestSizeBytes = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "HTTP request size in bytes",
		},
		[]string{"method", "path"},
	)

	httpResponseSizeBytes = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "HTTP response size in bytes",
		},
		[]string{"method", "path", "status"},
	)

	// 业务指标
	commandsExecutedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "commands_executed_total",
			Help: "Total number of commands executed",
		},
		[]string{"command", "status"},
	)

	queriesExecutedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "queries_executed_total",
			Help: "Total number of queries executed",
		},
		[]string{"query", "status"},
	)

	eventsPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_published_total",
			Help: "Total number of domain events published",
		},
		[]string{"event"},
	)

	// 数据库指标
	dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
		},
		[]string{"operation"},
	)
)

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

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		sr := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 记录请求大小
		if r.ContentLength > 0 {
			httpRequestSizeBytes.WithLabelValues(r.Method, path).Observe(float64(r.ContentLength))
		}

		next.ServeHTTP(sr, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(sr.statusCode)

		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		httpResponseSizeBytes.WithLabelValues(r.Method, path, status).Observe(float64(sr.bytesWritten))
	})
}

// RecordCommand 记录命令执行
func RecordCommand(command string, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	commandsExecutedTotal.WithLabelValues(command, status).Inc()
}

// RecordQuery 记录查询执行
func RecordQuery(query string, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	queriesExecutedTotal.WithLabelValues(query, status).Inc()
}

// RecordEvent 记录事件发布
func RecordEvent(event string) {
	eventsPublishedTotal.WithLabelValues(event).Inc()
}

// RecordDBQuery 记录数据库查询
func RecordDBQuery(operation string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	dbQueriesTotal.WithLabelValues(operation, status).Inc()
	dbQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}
