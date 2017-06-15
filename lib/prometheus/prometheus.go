package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var CountPayments = prometheus.NewCounter(prometheus.CounterOpts{Name: "billinger_count_payments", Help: "Количество пополнений баланса"})
var CountCreates = prometheus.NewCounter(prometheus.CounterOpts{Name: "billinger_count_creates", Help: "Количество генераций форм пополнения"})
var CountSelects = prometheus.NewCounter(prometheus.CounterOpts{Name: "billinger_count_selects", Help: "Количество получений статистики пополнений"})
var CountVerifyAuth = prometheus.NewCounter(prometheus.CounterOpts{Name: "billinger_count_verify_auth", Help: "Количество проверок сесии"})
var TarantoolResponseTime = prometheus.NewGauge(prometheus.GaugeOpts{Name: "billinger_tarantool_response_time", Help: "Время ответа от Tarantool"})

func init() {
	prometheus.MustRegister(CountPayments)
	prometheus.MustRegister(CountCreates)
	prometheus.MustRegister(CountSelects)
	prometheus.MustRegister(CountVerifyAuth)
	prometheus.MustRegister(TarantoolResponseTime)
}

func GetMetrics(c *gin.Context) {
	prom := prometheus.UninstrumentedHandler()
	prom.ServeHTTP(c.Writer, c.Request)
}
