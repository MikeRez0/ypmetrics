package handlers

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mh *MetricsHandler) MetricListView(c *gin.Context) {
	metrics := mh.Store.Metrics()

	err := mh.Template.ExecuteTemplate(c.Writer, "T", metrics)
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		log.Printf("failed to parse the metrics HTML template: %v", err)
		return
	}
}
