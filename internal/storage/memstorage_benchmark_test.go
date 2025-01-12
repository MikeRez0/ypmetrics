package storage_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func BenchmarkMemStorage(b *testing.B) {
	s := storage.NewMemStorage()

	for range b.N {
		for j := range 100 {
			_, err := s.UpdateCounter(context.Background(), "testCounter"+strconv.Itoa(j%10), 5)
			if err != nil {
				b.Fatal(err)
			}
			_, err = s.UpdateGauge(context.Background(), "testGauge"+strconv.Itoa(j%10), model.GaugeValue(j*5))
			if err != nil {
				b.Fatal(err)
			}
		}
		for j := 99; j >= 0; j -= 2 {
			_, err := s.GetCounter(context.Background(), "testCounter"+strconv.Itoa(j%10))
			if err != nil {
				b.Fatal(err)
			}
			_, err = s.GetGauge(context.Background(), "testGauge"+strconv.Itoa(j%10))
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
