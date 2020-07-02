package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"rampart/internal/config"
)

func newHistogramTracker(miner string, config *config.HistogramTracker) *histogramTracker {
	histogramVec := promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: config.Name, Help: config.Help, Buckets: config.Buckets},
		config.Labels,
	)
	return &histogramTracker{histogramVec, histogramVec.WithLabelValues(miner), config.Name, miner}
}

type histogramTracker struct {
	histogramVec *prometheus.HistogramVec
	observer     prometheus.Observer
	name         string
	miner        string
}

func (tracker *histogramTracker) track(value float64) {
	tracker.observer.Observe(value)
}

func (tracker *histogramTracker) unregister() error {
	if tracker.histogramVec.DeleteLabelValues(tracker.miner) &&
		prometheus.Unregister(tracker.histogramVec) {
		return nil
	}
	return fmt.Errorf("metrics: %s histogram tracker failed to unregister", tracker.name)
}
