package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PingDB - Handler for database ping.
func (mh *MetricsHandler) PingDB(c *gin.Context) {
	err := mh.Store.Ping()
	if err != nil {
		err := c.AbortWithError(500, err)
		mh.Log.Error("Ping DB error", zap.Error(err))
		return
	}
	c.Status(200)
}
