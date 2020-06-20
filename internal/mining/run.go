package mining

import (
	"fmt"
	"rampart/internal/database"
	"rampart/internal/mining/config"
	"rampart/internal/mining/domria"
	"rampart/internal/misc"
)

func Run() error {
	cfg, err := config.NewMining()
	if err != nil {
		return err
	}
	db, err := database.Setup(cfg.Params)
	if err != nil {
		return err
	}
	prospector := domria.NewProspector(misc.Secondary, cfg.Prospectors.Domria, db)
	if err = prospector.Prospect(); err != nil {
		_ = db.Close()
		return err
	}
	if err = db.Close(); err != nil {
		return fmt.Errorf("mining: failed to close the db, %v", err)
	}
	return nil
}
