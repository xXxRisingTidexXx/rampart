package config

import (
	"fmt"
)

type Gatherer struct {
	MiningDurationTracker *HistogramTracker `yaml:"miningDurationTracker"`
}

func (gatherer *Gatherer) String() string {
	return fmt.Sprintf("{%v}", gatherer.MiningDurationTracker)
}
