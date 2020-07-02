package metrics

import (
	"rampart/internal/config"
)

func NewGatherer(miner string, config *config.Gatherer) *Gatherer {
	return &Gatherer{newHistogramTracker(miner, config.MiningDurationTracker)}
}

type Gatherer struct {
	miningDurationTracker *histogramTracker
}

func (gatherer *Gatherer) GatherMiningDuration(miningDuration float64) {
	gatherer.miningDurationTracker.track(miningDuration)
}

func (gatherer *Gatherer) Unregister() error {
	return gatherer.miningDurationTracker.unregister()
}
