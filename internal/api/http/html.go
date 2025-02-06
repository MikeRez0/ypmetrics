package http

import (
	_ "embed"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MikeRez0/ypmetrics/internal/model"
)

// MetricListView - Handler for html-page, containing metrics list.
func (mh *MetricsHandler) MetricListView(c *gin.Context) {
	type NV struct {
		Name  string
		Value string
	}
	metrics := mh.service.Metrics()
	metricStrings := make([]NV, len(metrics))

	for i, m := range metrics {
		switch m.MType {
		case model.CounterType:
			metricStrings[i] = NV{m.ID, strconv.FormatInt(*m.Delta, 10)}
		case model.GaugeType:
			metricStrings[i] = NV{m.ID, strconv.FormatFloat(*m.Value, 'f', 2, 64)}
		}
	}

	// Обсудить на 1-1 эту строку. Если ее поставить после ExecuteTemplate, то все "ломается"
	c.Writer.Header().Set("Content-Type", "text/html")
	err := mh.Template.ExecuteTemplate(c.Writer, "T", metricStrings)
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		log.Printf("failed to parse the metrics HTML template: %v", err)
		return
	}
	c.Status(http.StatusOK)
}
