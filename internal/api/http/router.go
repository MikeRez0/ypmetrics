package http

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/utils/netctrl"
)

// SetupRouter - create gin router with handlers.
func SetupRouter(h *MetricsHandler, mylog *zap.Logger, ipControl *netctrl.IPControl) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.GinLogger(mylog))
	r.HandleMethodNotAllowed = true

	if ipControl != nil {
		r.Use(ipControl.Handler())
	}

	r.GET("/", gzip.Gzip(gzip.DefaultCompression), h.MetricListView)

	r.POST("/update/:metricType/:metric/:value", h.UpdateMetricPlain)
	r.GET("/value/:metricType/:metric", h.GetMetricPlain)

	jsonGroup := r.Group("/")
	jsonGroup.Use(GinCompress(logger.LoggerWithComponent(mylog, "compress")))
	jsonGroup.POST("/update/", h.UpdateMetricJSON)
	jsonGroup.POST("/value/", h.GetMetricJSON)
	jsonGroup.POST("/updates/", h.BatchUpdateMetricsJSON)

	r.GET("/ping", h.PingDB)

	pprof.Register(r)

	return r
}
