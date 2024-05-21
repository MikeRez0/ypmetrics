package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateCounter(t *testing.T) {
	ms := NewMemStorage()

	testMetric := "testCounter"
	ms.UpdateCounter(testMetric, 1)
	val, err := ms.GetCounter(testMetric)
	assert.NoError(t, err)
	assert.Equal(t, CounterValue(1), val)

	ms.UpdateCounter(testMetric, 5)
	val, err = ms.GetCounter(testMetric)
	assert.NoError(t, err)
	assert.Equal(t, CounterValue(6), val)

	_, err = ms.GetCounter(testMetric + "_fake")
	assert.Error(t, err)
}

func TestMemStorage_UpdateGauge(t *testing.T) {
	ms := NewMemStorage()

	testMetric := "testGauge"
	ms.UpdateGauge(testMetric, 1)
	val, err := ms.GetGauge(testMetric)
	assert.NoError(t, err)
	assert.Equal(t, GaugeValue(1), val)

	ms.UpdateGauge(testMetric, 5)
	val, err = ms.GetGauge(testMetric)
	assert.NoError(t, err)
	assert.Equal(t, GaugeValue(5), val)

	_, err = ms.GetGauge(testMetric + "_fake")
	assert.Error(t, err)
}
