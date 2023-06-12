package middlewares

import (
	"strconv"

	"github.com/JohnSalazar/microservices-go-common/metrics"
	"github.com/JohnSalazar/microservices-go-common/services"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Metrics(service services.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		appMetric := metrics.NewHttpMetrics(c.Request.URL.Path, c.Request.Method)
		appMetric.Started()
		c.Next()

		response := c.Writer
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(response.Status())
		service.SaveHttp(appMetric)
	}
}

func MetricsHandler() gin.HandlerFunc {
	handler := promhttp.Handler()

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
