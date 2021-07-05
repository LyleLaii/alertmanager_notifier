package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// GaugeVecAPIDuration api duration time, gauge type
	GaugeVecAPIDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiDuration",
		Help: "api耗时单位ms",
	}, []string{"WSorAPI"})
	// GaugeVecAPIMethod api method count, gauge type
	GaugeVecAPIMethod = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiCount",
		Help: "各种网络请求次数",
	}, []string{"method"})
	// GaugeVecAPIError api error count, gauge type
	GaugeVecAPIError = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiErrorCount",
		Help: "请求api错误的次数",
	}, []string{"function"})
	// CountVecNotifier notifier count, count type
	CountVecNotifier = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "notifyCount",
		Help: "告警发送次数",
	}, []string{"type"})
	// CountVecErrorNotifier notifier count, count type
	CountVecErrorNotifier = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "notifyErrorCount",
		Help: "告警发送异常次数",
	}, []string{"type"})
	// HistogramVecNofityDuration notify duration, histogram type
	HistogramVecNofityDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "nofityDuration",
		Help:    "告警发送耗时",
		Buckets: []float64{1000, 10000, 30000, 60000}, // 1s, 10s, 30s, 60s
	}, []string{"type"})
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(GaugeVecAPIMethod, GaugeVecAPIDuration, GaugeVecAPIError, CountVecNotifier, CountVecErrorNotifier, HistogramVecNofityDuration)
}


func Register(r gin.IRouter) {
	r.GET("metrics", gin.WrapH(promhttp.Handler()))
}