package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	cMetricNotFound         = "metric not found"
	cMetricTypeNotFound     = "metric type not found"
	cMetricTypeNameNotFound = "%s not a metric type"
)

func (mh *MetricsHandler) UpdateMetricPlain(c *gin.Context) {
	var (
		metricType = c.Param("metricType")
		metric     = c.Param("metric")
		valueRaw   = c.Param("value")
	)

	if metric == "" {
		// c.AbortWithStatus(http.StatusNotFound)
		handleError(c, http.StatusNotFound, errors.New(cMetricNotFound), mh.Log, cMetricNotFound)
		return
	}

	switch metricType {
	case model.GaugeType:
		value, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			// c.AbortWithStatus(http.StatusBadRequest)
			handleError(c, http.StatusBadRequest, err, mh.Log, "bad request")
			return
		}
		_, err = mh.Store.UpdateGauge(c, metric, model.GaugeValue(value))
		if err != nil {
			// err = c.AbortWithError(http.StatusInternalServerError, err)
			// mh.Log.Error("error on guage update", zap.Error(err))
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on gauge update")
			return
		}
	case model.CounterType:
		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			// c.AbortWithStatus(http.StatusBadRequest)
			handleError(c, http.StatusBadRequest, err, mh.Log, "bad request")
			return
		}
		_, err = mh.Store.UpdateCounter(c, metric, model.CounterValue(value))
		if err != nil {
			// err = c.AbortWithError(http.StatusInternalServerError, err)
			// mh.Log.Error("error on counter update", zap.Error(err))
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on counter update")
			return
		}
	default:
		// c.AbortWithStatus(http.StatusBadRequest)
		handleError(c, http.StatusBadRequest, fmt.Errorf(cMetricTypeNameNotFound, metricType), mh.Log, cMetricTypeNotFound)
		return
	}

	c.Status(http.StatusOK)
}

func (mh *MetricsHandler) GetMetricPlain(c *gin.Context) {
	var (
		metricType = c.Param("metricType")
		metric     = c.Param("metric")
	)

	switch metricType {
	case model.GaugeType:
		value, err := mh.Store.GetGauge(c, metric)
		if err != nil {
			// err = c.AbortWithError(http.StatusNotFound, err)
			// mh.Log.Error("error on get metric", zap.Error(err))
			handleError(c, http.StatusNotFound, err, mh.Log, cMetricNotFound)
			return
		}
		_, err = c.Writer.WriteString(strconv.FormatFloat(float64(value), 'f', -1, 64))
		if err != nil {
			// err = c.AbortWithError(http.StatusInternalServerError, err)
			// mh.Log.Error("error on get metric", zap.Error(err))
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on get metric")
			return
		}
	case model.CounterType:
		value, err := mh.Store.GetCounter(c, metric)
		if err != nil {
			// err = c.AbortWithError(http.StatusNotFound, err)
			// log.Println(err)
			handleError(c, http.StatusNotFound, err, mh.Log, cMetricNotFound)
			return
		}
		_, err = c.Writer.WriteString(strconv.FormatInt(int64(value), 10))
		if err != nil {
			// err = c.AbortWithError(http.StatusInternalServerError, err)
			// mh.Log.Error("error on get metric", zap.Error(err))
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on get metric")
			return
		}
	default:
		// err := c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s not a metric type", metricType))
		// if err != nil {
		// 	mh.Log.Error("metric type not found", zap.Error(err))
		// 	return
		// }
		handleError(c, http.StatusBadRequest, fmt.Errorf(cMetricTypeNameNotFound, metricType),
			mh.Log, cMetricTypeNotFound)
		return
	}

	c.Header("Content-Type", "text/plain")

	c.Status(http.StatusOK)
}

func (mh *MetricsHandler) UpdateMetricJSON(c *gin.Context) {
	var metric model.Metrics
	if err := c.ShouldBindJSON(&metric); err != nil {
		// _ = c.Error(err)
		// c.AbortWithStatus(http.StatusBadRequest)
		handleError(c, http.StatusBadRequest, err, mh.Log, "Bad request")
		return
	}

	if metric.ID == "" {
		// c.AbortWithStatus(http.StatusNotFound)
		handleError(c, http.StatusNotFound, errors.New(cMetricNotFound), mh.Log, cMetricNotFound)
	}

	switch metric.MType {
	case model.GaugeType:
		v, err := mh.Store.UpdateGauge(c, metric.ID, model.GaugeValue(*metric.Value))
		if err != nil {
			// err = c.AbortWithError(http.StatusInternalServerError, err)
			// mh.Log.Error("Error updating metric "+metric.ID, zap.Error(err))
			handleError(c, http.StatusInternalServerError, err, mh.Log, "Error updating metric: "+metric.ID)
			return
		}
		var newVal = float64(v)
		metric.Value = &newVal
	case model.CounterType:
		v, err := mh.Store.UpdateCounter(c, metric.ID, model.CounterValue(*metric.Delta))
		if err != nil {
			// err = c.AbortWithError(http.StatusInternalServerError, err)
			// mh.Log.Error("Error updating metric ", zap.String("MetricID", metric.ID), zap.Error(err))
			handleError(c, http.StatusInternalServerError, err, mh.Log, "Error updating metric: "+metric.ID)
			return
		}
		var newVal = int64(v)
		metric.Delta = &newVal
	default:
		// c.AbortWithStatus(http.StatusBadRequest)
		handleError(c, http.StatusBadRequest, fmt.Errorf(cMetricTypeNameNotFound, metric.MType), mh.Log, cMetricTypeNotFound)
		return
	}

	if mh.Signer != nil {
		h, err := mh.Signer.GetHashJSON(metric)
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "Error calculating hash")
			return
		}
		c.Header("HashSHA256", h)
	}

	c.JSON(http.StatusOK, metric)
}

func (mh *MetricsHandler) GetMetricJSON(c *gin.Context) {
	var metric model.Metrics
	if err := c.ShouldBindJSON(&metric); err != nil {
		handleError(c, http.StatusBadRequest, err, mh.Log, "Bad request")
		return
	}

	if metric.ID == "" {
		handleError(c, http.StatusNotFound, errors.New("metric not found"), mh.Log, "error")
	}

	switch metric.MType {
	case model.GaugeType:
		value, err := mh.Store.GetGauge(c, metric.ID)
		if err != nil {
			handleError(c, http.StatusNotFound, err, mh.Log, "error")
			return
		}
		metric.Value = (*float64)(&value)
	case model.CounterType:
		value, err := mh.Store.GetCounter(c, metric.ID)
		if err != nil {
			handleError(c, http.StatusNotFound, err, mh.Log, "error")
			return
		}
		metric.Delta = (*int64)(&value)
	default:
		handleError(c, http.StatusBadRequest, fmt.Errorf(cMetricTypeNameNotFound, metric.MType), mh.Log, "")
		return
	}

	if mh.Signer != nil {
		h, err := mh.Signer.GetHashJSON(metric)
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "Error calculating hash")
			return
		}
		c.Header("HashSHA256", h)
	}
	c.JSON(http.StatusOK, metric)
}

func (mh *MetricsHandler) BatchUpdateMetricsJSON(c *gin.Context) {
	var metrics []model.Metrics
	if err := c.ShouldBindJSON(&metrics); err != nil {
		// _ = c.Error(err)
		// c.AbortWithStatus(http.StatusBadRequest)
		handleError(c, http.StatusBadRequest, err, nil, "")
		return
	}

	hashReq := c.Request.Header.Get(model.HeaderSignerHash)
	if mh.Signer != nil && hashReq == "" {
		handleError(c, http.StatusBadRequest, errors.New("hash expected in header"), mh.Log, "")
		return
	}
	if hashReq != "" {
		mh.Log.Debug("Request hash", zap.String("Hash", hashReq))
	}
	if mh.Signer != nil && !mh.Signer.ValidateJSON(metrics, hashReq) {
		handleError(c, http.StatusBadRequest, errors.New("hash value error"), mh.Log, "")
		return
	}

	err := mh.Store.BatchUpdate(c, metrics)
	if err != nil {
		if errors.As(err, &model.BadValueError{}) {
			handleError(c, http.StatusBadRequest, err, mh.Log, "")
			return
		}
		handleError(c, http.StatusInternalServerError, err, mh.Log, "Error batch updating metrics")
		return
	}

	if mh.Signer != nil {
		h, err := mh.Signer.GetHashJSON(metrics)
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "Error calculating hash")
			return
		}
		c.Header(model.HeaderSignerHash, h)
	}

	c.JSON(http.StatusOK, metrics)
}

func handleError(c *gin.Context, statusCode int, err error, logger *zap.Logger, message string) {
	err = c.AbortWithError(statusCode, err)
	if logger != nil {
		if message != "" {
			logger.Error(message, zap.Error(err))
		} else {
			logger.Error("error (resp status code: "+strconv.Itoa(statusCode)+")", zap.Error(err))
		}

	}
}
