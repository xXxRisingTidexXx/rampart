package metrics

import (
	"rampart/internal/config"
)

func NewGatherer(miner string, config *config.Gatherer) *Gatherer {
	return &Gatherer{newTracker(miner)}
}

type Gatherer struct {
	miningDurationTracker *tracker
}

func (gatherer *Gatherer) GatherMiningDuration(miningDuration float64) {
	gatherer.miningDurationTracker.track(miningDuration)
}

func (gatherer *Gatherer) Unregister() error {
	return gatherer.miningDurationTracker.unregister()
}
