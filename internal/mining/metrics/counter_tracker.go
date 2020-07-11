package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"rampart/internal/config"
)

func newCounterTracker(config *config.CounterTracker, miner string, targets []string) *counterTracker {
	counterVec := promauto.NewCounterVec(
		prometheus.CounterOpts{Name: config.Name, Help: config.Help},
		config.Labels,
	)
	counterMap := make(map[string]prometheus.Counter, len(targets))
	for _, target := range targets {
		counterMap[target] = counterVec.WithLabelValues(miner, target)
	}
	return &counterTracker{counterVec, counterMap, config.Name, miner}
}

type counterTracker struct {
	counterVec *prometheus.CounterVec
	counterMap map[string]prometheus.Counter
	name       string
	miner      string
}

func (tracker *counterTracker) track(target string) {
	counter, ok := tracker.counterMap[target]
	if !ok {
		panic(fmt.Sprintf("metrics: counter tracker %s failed to track target %s", tracker.name, target))
	}
	counter.Inc()
}

func (tracker *counterTracker) unregister() error {
	err := fmt.Errorf("metrics: counter tracker %s failed to unregister", tracker.name)
	for target := range tracker.counterMap {
		if !tracker.counterVec.DeleteLabelValues(tracker.miner, target) {
			return err
		}
	}
	if prometheus.Unregister(tracker.counterVec) {
		return nil
	}
	return err
}
