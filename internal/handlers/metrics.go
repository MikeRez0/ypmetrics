package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
)

const (
	cMetricNotFound         = "metric not found"
	cMetricTypeNotFound     = "metric type not found"
	cMetricTypeNameNotFound = "%s not a metric type"
)

// UpdateMetricPlain - Update metric by plain text request.
func (mh *MetricsHandler) UpdateMetricPlain(c *gin.Context) {
	var (
		metricType = c.Param("metricType")
		metric     = c.Param("metric")
		valueRaw   = c.Param("value")
	)

	if metric == "" {
		handleError(c, http.StatusNotFound, errors.New(cMetricNotFound), mh.Log, cMetricNotFound)
		return
	}

	switch metricType {
	case model.GaugeType:
		value, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			handleError(c, http.StatusBadRequest, err, mh.Log, "bad request")
			return
		}
		_, err = mh.Store.UpdateGauge(c, metric, model.GaugeValue(value))
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on gauge update")
			return
		}
	case model.CounterType:
		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			handleError(c, http.StatusBadRequest, err, mh.Log, "bad request")
			return
		}
		_, err = mh.Store.UpdateCounter(c, metric, model.CounterValue(value))
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on counter update")
			return
		}
	default:
		handleError(c, http.StatusBadRequest, fmt.Errorf(cMetricTypeNameNotFound, metricType), mh.Log, cMetricTypeNotFound)
		return
	}

	c.Status(http.StatusOK)
}

// GetMetricPlain - Get metric by plain text request.
func (mh *MetricsHandler) GetMetricPlain(c *gin.Context) {
	var (
		metricType = c.Param("metricType")
		metric     = c.Param("metric")
	)

	switch metricType {
	case model.GaugeType:
		value, err := mh.Store.GetGauge(c, metric)
		if err != nil {
			handleError(c, http.StatusNotFound, err, mh.Log, cMetricNotFound)
			return
		}
		_, err = c.Writer.WriteString(strconv.FormatFloat(float64(value), 'f', -1, 64))
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on get metric")
			return
		}
	case model.CounterType:
		value, err := mh.Store.GetCounter(c, metric)
		if err != nil {
			handleError(c, http.StatusNotFound, err, mh.Log, cMetricNotFound)
			return
		}
		_, err = c.Writer.WriteString(strconv.FormatInt(int64(value), 10))
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "error on get metric")
			return
		}
	default:
		handleError(c, http.StatusBadRequest, fmt.Errorf(cMetricTypeNameNotFound, metricType),
			mh.Log, cMetricTypeNotFound)
		return
	}

	c.Header("Content-Type", "text/plain")

	c.Status(http.StatusOK)
}

// UpdateMetricJSON - Update metric by JSON request.
func (mh *MetricsHandler) UpdateMetricJSON(c *gin.Context) {
	var metric model.Metrics
	if err := c.ShouldBindJSON(&metric); err != nil {
		handleError(c, http.StatusBadRequest, err, mh.Log, "Bad request")
		return
	}

	if metric.ID == "" {
		handleError(c, http.StatusNotFound, errors.New(cMetricNotFound), mh.Log, cMetricNotFound)
	}

	switch metric.MType {
	case model.GaugeType:
		v, err := mh.Store.UpdateGauge(c, metric.ID, model.GaugeValue(*metric.Value))
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "Error updating metric: "+metric.ID)
			return
		}
		var newVal = float64(v)
		metric.Value = &newVal
	case model.CounterType:
		v, err := mh.Store.UpdateCounter(c, metric.ID, model.CounterValue(*metric.Delta))
		if err != nil {
			handleError(c, http.StatusInternalServerError, err, mh.Log, "Error updating metric: "+metric.ID)
			return
		}
		var newVal = int64(v)
		metric.Delta = &newVal
	default:
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

// GetMetricJSON - Get metric by JSON request.
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

// BatchUpdateMetricsJSON - Update multiple metrics by JSON request.
func (mh *MetricsHandler) BatchUpdateMetricsJSON(c *gin.Context) {
	var metrics []model.Metrics

	hashReq := c.Request.Header.Get(model.HeaderSignerHash)
	if mh.Signer != nil && hashReq == "" {
		handleError(c, http.StatusBadRequest, errors.New("hash expected in header"), mh.Log, "")
		return
	}

	encryptKey := c.Request.Header.Get(model.HeaderEncryptKey)
	if mh.Decrypter != nil && encryptKey == "" {
		handleError(c, http.StatusBadRequest, errors.New("encrypt key expected in header"), mh.Log, "")
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		handleError(c, http.StatusBadRequest, err, mh.Log, "")
		return
	}

	if mh.Decrypter != nil {
		key, err := base64.StdEncoding.DecodeString(encryptKey)
		if err != nil {
			handleError(c, http.StatusBadRequest, err, mh.Log, "")
		}
		encData, err := base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			handleError(c, http.StatusBadRequest, err, mh.Log, "")
		}
		d, err := mh.Decrypter.Decrypt(&signer.Envelope{
			Key:  key,
			Data: encData,
		})
		if err != nil {
			handleError(c, http.StatusBadRequest, err, mh.Log, "")
		}
		data = d
	}

	if mh.Signer != nil && !mh.Signer.Validate(data, hashReq) {
		handleError(c, http.StatusBadRequest, errors.New("hash value error"), mh.Log, "")
		return
	}

	if err = json.Unmarshal(data, &metrics); err != nil {
		// if err = c.ShouldBindJSON(&metrics); err != nil {
		handleError(c, http.StatusBadRequest, err, nil, "")
		return
	}

	err = mh.Store.BatchUpdate(c, metrics)
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

// handleError - helper for error handle.
func handleError(c *gin.Context, statusCode int, err error, logger *zap.Logger, message string) {
	err = c.AbortWithError(statusCode, err)
	if logger != nil {
		msg := message
		if message == "" {
			msg = "error (resp status code: " + strconv.Itoa(statusCode) + ")"
		}
		if statusCode == http.StatusInternalServerError {
			logger.Error(msg, zap.Error(err))
		} else {
			logger.Debug(msg, zap.Error(err))
		}
	}
}
