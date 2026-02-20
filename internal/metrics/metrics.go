package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: 		"http_request_duration_seconds",
			Help:    	"HTTP request latency in seconds",
			Buckets: 	prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// Cache metrics
    CacheHits = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_hits_total",
            Help: "Total number of cache hits",
        },
        []string{"cache_key"},
    )
    
    CacheMisses = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_misses_total",
            Help: "Total number of cache misses",
        },
        []string{"cache_key"},
    )
    
    // Database metrics
    DbQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_queries_total",
            Help: "Total number of database queries",
        },
        []string{"operation"},
    )
    
    DbQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query latency in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"operation"},
    )
    
    // Product metrics
    ProductsCreated = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "products_created_total",
            Help: "Total number of products created",
        },
    )
    
    ProductsDeleted = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "products_deleted_total",
            Help: "Total number of products deleted",
        },
    )
    
    // Auth metrics
    UserRegistrations = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "user_registrations_total",
            Help: "Total number of user registrations",
        },
    )
    
    LoginAttempts = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "login_attempts_total",
            Help: "Total number of login attempts",
        },
        []string{"status"}, // success or failure
    )
)

func TimeDatabaseQuery(operation string) func() {
    start := time.Now()
    return func() {
        duration := time.Since(start).Seconds()
        DbQueryDuration.WithLabelValues(operation).Observe(duration)
        DbQueriesTotal.WithLabelValues(operation).Inc()
    }
}