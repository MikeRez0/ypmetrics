package storage

import (
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateCounter(t *testing.T) { //nolint:dupl //false positive
	ms := NewMemStorage()

	const testMetricCounter = "testCounter"
	val, err := ms.UpdateCounter(testMetricCounter, 1)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(1), val)
	val, err = ms.GetCounter(testMetricCounter)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(1), val)

	val, err = ms.UpdateCounter(testMetricCounter, 5)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(6), val)
	val, err = ms.GetCounter(testMetricCounter)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(6), val)

	_, err = ms.GetCounter(testMetricCounter + "_fake")
	assert.Error(t, err)
}

func TestMemStorage_UpdateGauge(t *testing.T) { //nolint:dupl //false positive
	ms := NewMemStorage()

	const testMetricGauge = "testGauge"
	val, err := ms.UpdateGauge(testMetricGauge, 1)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(1), val)

	val, err = ms.GetGauge(testMetricGauge)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(1), val)

	val, err = ms.UpdateGauge(testMetricGauge, 5)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(1), val)

	val, err = ms.GetGauge(testMetricGauge)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(5), val)

	_, err = ms.GetGauge(testMetricGauge + "_fake")
	assert.Error(t, err)
}
