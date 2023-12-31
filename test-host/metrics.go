package main

import (
    "net/http"
    "strconv"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Number of requests.",
    },
    []string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "response_status",
        Help: "Status of HTTP response",
    },
    []string{"status"},
)

var httpDuration = promauto.NewHistogramVec(
    prometheus.HistogramOpts{
        Name: "http_response_time_seconds",
        Help: "Duration of HTTP requests.",
    },
    []string{"path"},
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
    return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func prometheusMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path

        timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
        rw := NewResponseWriter(w)
        next.ServeHTTP(rw, r)

        statusCode := rw.statusCode

        totalRequests.WithLabelValues(path).Inc()
        responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
        totalRequests.WithLabelValues(path).Inc()

        timer.ObserveDuration()
    })
}
