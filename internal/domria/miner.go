package domria

import (
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining"
)

func NewMiner(config config.DomriaMiner) mining.Miner {
	return &miner{config.Name, config.Spec}
}

type miner struct {
	name string
	spec string
}

func (m *miner) Name() string {
	return m.name
}

func (m *miner) Spec() string {
	return m.spec
}

func (m *miner) MineFlat() (mining.Flat, error) {
	return mining.Flat{}, nil
}
