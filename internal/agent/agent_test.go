package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRuntimeMetrics(t *testing.T) {
	ms := NewMetricStore()
	ReadRuntimeMetrics(ms)
	for _, v := range runtimeMetricNames {

		assert.Contains(t, ms.Metrics, v)
	}
}
