package middleware

import (
	stat "alertmanager_notifier/metrics"

	"time"

	"github.com/gin-gonic/gin"
)

// MwPrometheusHTTP prometheus http handler
func MwPrometheusHTTP(c *gin.Context) {
	start := time.Now()
	method := c.Request.Method
	stat.GaugeVecAPIMethod.WithLabelValues(method).Inc()

	c.Next()
	// after request
	end := time.Now()
	d := end.Sub(start) / time.Millisecond
	stat.GaugeVecAPIDuration.WithLabelValues(method).Set(float64(d))
}
