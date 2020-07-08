package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"rampart/internal/config"
)

func newHistogramTracker(
	config *config.HistogramTracker,
	miner string,
	targets []string,
) *histogramTracker {
	histogramVec := promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: config.Name, Help: config.Help, Buckets: config.Buckets},
		config.Labels,
	)
	observerMap := make(map[string]prometheus.Observer, len(targets))
	for _, target := range targets {
		observerMap[target] = histogramVec.WithLabelValues(miner, target)
	}
	return &histogramTracker{histogramVec, observerMap, config.Name, miner}
}

type histogramTracker struct {
	histogramVec *prometheus.HistogramVec
	observerMap  map[string]prometheus.Observer
	name         string
	miner        string
}

func (tracker *histogramTracker) track(target string, value float64) {
	observer, ok := tracker.observerMap[target]
	if !ok {
		panic(
			fmt.Sprintf("metrics: %s histogram tracker failed to track target %s", tracker.name, target),
		)
	}
	observer.Observe(value)
}

func (tracker *histogramTracker) unregister() error {
	err := fmt.Errorf("metrics: histogram tracker %s failed to unregister", tracker.name)
	for target := range tracker.observerMap {
		if !tracker.histogramVec.DeleteLabelValues(tracker.miner, target) {
			return err
		}
	}
	if prometheus.Unregister(tracker.histogramVec) {
		return nil
	}
	return err
}
