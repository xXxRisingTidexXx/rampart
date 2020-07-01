package config

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type exponentialBuckets struct {
	Start  float64 `yaml:"start"`
	Factor float64 `yaml:"factor"`
	Count  int     `yaml:"count"`
}

func (buckets exponentialBuckets) toFloats() (floats []float64, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			floats = nil
			err = fmt.Errorf("config: exponential buckets failed to convert to floats, %v", err1)
		}
	}()
	return prometheus.ExponentialBuckets(buckets.Start, buckets.Factor, buckets.Count), nil
}

func (buckets exponentialBuckets) String() string {
	return fmt.Sprintf("{%.3f %.3f %d}", buckets.Start, buckets.Factor, buckets.Count)
}
