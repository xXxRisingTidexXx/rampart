package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func newTracker(miner string) *tracker {
	histogramVec := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "rampart_mining_duration",
			Help: "Total elapsed time of a single miner execution in seconds",
		},
		[]string{"miner"},
	)
	return &tracker{histogramVec, histogramVec.WithLabelValues(miner), "rampart_mining_duration", miner}
}

type tracker struct {
	histogramVec *prometheus.HistogramVec
	observer     prometheus.Observer
	name         string
	miner        string
}

func (tracker *tracker) track(value float64) {
	tracker.observer.Observe(value)
}

func (tracker *tracker) unregister() error {
	if tracker.histogramVec.DeleteLabelValues(tracker.miner) &&
		prometheus.Unregister(tracker.histogramVec) {
		return nil
	}
	return fmt.Errorf("metrics: %s tracker failed to unregister", tracker.name)
}
