package config

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type linearBuckets struct {
	Start float64 `yaml:"start"`
	Width float64 `yaml:"width"`
	Count int     `yaml:"count"`
}

func (buckets linearBuckets) toFloats() (floats []float64, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			floats = nil
			err = fmt.Errorf("config: linear buckets failed to convert to floats, %v", err1)
		}
	}()
	return prometheus.LinearBuckets(buckets.Start, buckets.Width, buckets.Count), nil
}

func (buckets linearBuckets) String() string {
	return fmt.Sprintf("{%.3f %.3f %d}", buckets.Start, buckets.Width, buckets.Count)
}
