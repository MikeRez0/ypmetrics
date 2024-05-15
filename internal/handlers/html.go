package handlers

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed "templates/metrics.html"
var templateContent string
var tmpl *template.Template

func (mh *MetricsHandler) MetricListView(c *gin.Context) {
	metrics := mh.Store.Metrics()

	if tmpl == nil {
		var err error
		tmpl, err = template.New("metrics").Parse(templateContent)
		if err != nil {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}
	}
	err := tmpl.ExecuteTemplate(c.Writer, "T", metrics)
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
}
