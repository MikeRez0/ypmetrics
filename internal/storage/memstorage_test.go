package storage

import (
	"context"
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateCounter(t *testing.T) { //nolint:dupl //false positive
	ms := NewMemStorage()

	const testMetricCounter = "testCounter"
	val, err := ms.UpdateCounter(context.Background(), testMetricCounter, 1)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(1), val)
	val, err = ms.GetCounter(context.Background(), testMetricCounter)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(1), val)

	val, err = ms.UpdateCounter(context.Background(), testMetricCounter, 5)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(6), val)
	val, err = ms.GetCounter(context.Background(), testMetricCounter)
	assert.NoError(t, err)
	assert.Equal(t, model.CounterValue(6), val)

	_, err = ms.GetCounter(context.Background(), testMetricCounter+"_fake")
	assert.Error(t, err)
}

func TestMemStorage_UpdateGauge(t *testing.T) { //nolint:dupl //false positive
	ms := NewMemStorage()

	const testMetricGauge = "testGauge"
	val, err := ms.UpdateGauge(context.Background(), testMetricGauge, 1)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(1), val)

	val, err = ms.GetGauge(context.Background(), testMetricGauge)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(1), val)

	val, err = ms.UpdateGauge(context.Background(), testMetricGauge, 5)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(5), val)

	val, err = ms.GetGauge(context.Background(), testMetricGauge)
	assert.NoError(t, err)
	assert.Equal(t, model.GaugeValue(5), val)

	_, err = ms.GetGauge(context.Background(), testMetricGauge+"_fake")
	assert.Error(t, err)
}
