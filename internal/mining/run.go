package mining

import (
	"fmt"
	"rampart/internal/config"
	"rampart/internal/database"
	"rampart/internal/mining/domria"
	"rampart/internal/misc"
	"rampart/internal/secrets"
)

func Run() error {
	scr, err := secrets.NewSecrets()
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	db, err := database.Setup(scr.DSN, cfg.Mining.DSNParams)
	if err != nil {
		return err
	}
	prospector := domria.NewProspector(misc.Secondary, cfg.Mining.Prospectors.Domria, db)
	if err = prospector.Prospect(); err != nil {
		_ = db.Close()
		return err
	}
	if err = db.Close(); err != nil {
		return fmt.Errorf("mining: failed to close the db, %v", err)
	}
	return nil
}
