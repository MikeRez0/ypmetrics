package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/gin-gonic/gin"
)

func (mh *MetricsHandler) UpdateMetricGin(c *gin.Context) {
	var (
		metricType = c.Param("metricType")
		metric     = c.Param("metric")
		valueRaw   = c.Param("value")
	)

	if metric == "" {
		c.AbortWithStatus(http.StatusNotFound)
	}

	switch metricType {
	case storage.GaugeType:
		value, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		mh.Store.UpdateGauge(metric, storage.GaugeValue(value))
	case storage.CounterType:
		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		mh.Store.UpdateCounter(metric, storage.CounterValue(value))
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (mh *MetricsHandler) GetMetricGin(c *gin.Context) {
	var (
		metricType = c.Param("metricType")
		metric     = c.Param("metric")
	)

	switch metricType {
	case storage.GaugeType:
		value, err := mh.Store.GetGauge(metric)
		if err != nil {
			err = c.AbortWithError(http.StatusNotFound, err)
			log.Println(err)
			return
		}
		_, err = c.Writer.WriteString(strconv.FormatFloat(float64(value), 'f', -1, 64))
		if err != nil {
			log.Println(err)
			return
		}
	case storage.CounterType:
		value, err := mh.Store.GetCounter(metric)
		if err != nil {
			err = c.AbortWithError(http.StatusNotFound, err)
			log.Println(err)
			return
		}
		_, err = c.Writer.WriteString(strconv.FormatInt(int64(value), 10))
		if err != nil {
			log.Println(err)
			return
		}
	default:
		err := c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s not a metric type", metricType))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	c.Header("Content-Type", "text/plain")

	c.Status(http.StatusOK)
}
