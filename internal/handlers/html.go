package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mh *MetricsHandler) MetricListView(c *gin.Context) {
	metrics := mh.Store.Metrics()

	var tmplFile = "../../internal/templates/metrics.html"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = tmpl.ExecuteTemplate(c.Writer, "T", metrics)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
