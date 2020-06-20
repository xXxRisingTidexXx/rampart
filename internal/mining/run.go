package mining

import (
	"fmt"
	"rampart/internal/config"
	"rampart/internal/database"
	"rampart/internal/mining/domria"
	"rampart/internal/misc"
)

func Run() error {
	rampart, err := config.NewRampart()
	if err != nil {
		return err
	}
	db, err := database.Setup(rampart.Mining.DSNParams)
	if err != nil {
		return err
	}
	prospector := domria.NewProspector(misc.Secondary, rampart.Mining.Prospectors.Domria, db)
	if err = prospector.Prospect(); err != nil {
		_ = db.Close()
		return err
	}
	if err = db.Close(); err != nil {
		return fmt.Errorf("mining: failed to close the db, %v", err)
	}
	return nil
}
