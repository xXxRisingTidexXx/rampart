package mining

import (
	"rampart/internal/mining/config"
	"rampart/internal/mining/domria"
	"rampart/internal/mining/misc"
)

func Run() error {
	cfg, err := config.NewMining()
	if err != nil {
		return err
	}
	prospector := domria.NewProspector(misc.Secondary, cfg.Prospectors.Domria)
	if err = prospector.Prospect(); err != nil {
		return err
	}
	return nil
}
