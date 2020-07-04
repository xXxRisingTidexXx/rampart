package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"rampart/internal/config"
)

func newCounterTracker(miner string, config *config.CounterTracker) *counterTracker {
	counterVec := promauto.NewCounterVec(
		prometheus.CounterOpts{Name: config.Name, Help: config.Help},
		config.Labels,
	)
	return &counterTracker{counterVec, counterVec.WithLabelValues(miner), config.Name, miner}
}

type counterTracker struct {
	counterVec *prometheus.CounterVec
	counter    prometheus.Counter
	name       string
	miner      string
}

func (tracker *counterTracker) track(value float64) {
	tracker.counter.Add(value)
}

func (tracker *counterTracker) unregister() error {
	if tracker.counterVec.DeleteLabelValues(tracker.miner) &&
		prometheus.Unregister(tracker.counterVec) {
		return nil
	}
	return fmt.Errorf("metrics: %s counter tracker failed to unregister", tracker.name)
}
