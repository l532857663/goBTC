package prometheus

import (
	"goBTC/global"
	"goBTC/models"
	"goBTC/utils/logutils"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 指标类型
/* Counter 计数器
描述：一个单调递增的计数器，仅能增加或重置为零，用于计数事件的发生次数。
典型应用：请求次数、任务完成次数、错误发生次数等。
*/
// RequestsTotal = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "http_requests_total",
// 		Help: "Total number of HTTP requests processed, labeled by status code and method.",
// 	},
// 	[]string{"code", "method"},
// )

/* Gauge 标尺
描述：一个表示单个数值的指标，可以任意增减，用于测量瞬时值。
典型应用：温度、内存使用、并发连接数等。
*/
// MemoryUsage = prometheus.NewGauge(
// 	prometheus.GaugeOpts{
// 		Name: "memory_usage_bytes",
// 		Help: "Current memory usage in bytes",
//		ConstLabels: prometheus.Labels{ // 常量标签
//	        "host": "localhost",
//	    },
// 	},
// )

/* Histogram 直方图
描述：一个累积直方图，用于记录观察值，并将其划分到不同的桶中，适合于观测事件的分布。
典型应用：请求延迟、响应时间等。
*/
// RequestDurationHistogram = prometheus.NewHistogramVec(
// 	prometheus.HistogramOpts{
// 		Name:    "http_request_duration_seconds",
// 		Help:    "Histogram of response latency (seconds) of HTTP requests.",
// 		Buckets: prometheus.DefBuckets,
// 	},
// 	[]string{"handler"},
// )

/* Summary 摘要
描述：类似于 Histogram，但更侧重于计算分位数和总数，适合于精确观测事件的分布。
典型应用：响应时间的分位数（如P50、P90、P99）、请求大小等。
*/
// requestDurationSummary = prometheus.NewSummary(
//     prometheus.SummaryOpts{
//         Name: "http_request_duration_seconds",
//         Help: "Duration of HTTP requests in seconds",
//     },
// )

var (
	// goroutine启动数量监控
	GoroutineUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "server_goroutine_usage_rate",
			Help: "Frequency of using goroutime by the service",
			ConstLabels: prometheus.Labels{
				"test_server": "使用的服务",
			},
		},
	)
)

func init() {
	// Register metrics with Prometheus's default registry.
	prometheus.MustRegister(GoroutineUsage)
}

func InitService(conf models.ServiceConf) {
	http.Handle("/metrics", promhttp.Handler())
	logutils.LogInfof(global.LOG, "InitService prometheus: %+v", conf.ServicePrometheusAddr)
	http.ListenAndServe(conf.ServicePrometheusAddr, nil)
}
