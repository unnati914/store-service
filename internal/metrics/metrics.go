package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "HTTP requests"},
		[]string{"route", "code"},
	)
)

func MustRegisterAll() {
	prometheus.MustRegister(HTTPRequests)
}
