package handlers

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mh *MetricsHandler) MetricListView(c *gin.Context) {
	metrics := mh.Store.MetricStrings()

	// Обсудить на 1-1 эту строку. Если ее поставить после ExecuteTemplate, то все "ломается"
	c.Writer.Header().Set("Content-Type", "text/html")
	err := mh.Template.ExecuteTemplate(c.Writer, "T", metrics)
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		log.Printf("failed to parse the metrics HTML template: %v", err)
		return
	}
	c.Status(http.StatusOK)
}
