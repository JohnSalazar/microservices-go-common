package services

import (
	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	SaveHttp(http *metrics.HttpMetrics)
	// SaveClient(client *metrics.ClientMetrics) error
}

type metricsService struct {
	config               *config.Config
	pHistogram           *prometheus.HistogramVec
	httpRequestHistogram *prometheus.HistogramVec
}

func NewMetricsService(
	config *config.Config,
) (*metricsService, error) {
	client := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.AppName + "_pushgateway",
		Subsystem: config.AppName,
		Name:      "cmd_duration_seconds",
		Help:      "Client application execution in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"name"})

	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.AppName + "_http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	service := &metricsService{
		config:               config,
		pHistogram:           client,
		httpRequestHistogram: http,
	}

	err := prometheus.Register(service.pHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	err = prometheus.Register(service.httpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	return service, nil
}

func (service *metricsService) SaveHttp(http *metrics.HttpMetrics) {
	service.httpRequestHistogram.WithLabelValues(http.Handler, http.Method, http.StatusCode).Observe(http.Duration)
}

// func (service *metricsService) SaveClient(client *metrics.ClientMetrics) error {
// 	gatewayURL := service.config.Prometheus.PROMETHEUS_PUSHGATEWAY
// 	service.pHistogram.WithLabelValues(client.Name).Observe(client.Duration)

// 	return push.New(gatewayURL, "cmd_job").Collector(service.pHistogram).Push()
// }
