package mining

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"rampart/internal/config"
	"rampart/internal/database"
	"rampart/internal/mining/domria"
	"rampart/internal/misc"
	"rampart/internal/secrets"
)

func Schedule() error {
	scr, err := secrets.NewSecrets()
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	db, err := database.NewDatabase(scr.DSN, cfg.Mining.DSNParams)
	if err != nil {
		return err
	}
	scheduler := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	_, err = scheduler.AddJob("* * * * *", domria.NewRunner(misc.Primary, cfg.Mining.Runners.Domria, db))
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("mining: failed to run domria primary runner, %v", err)
	}
	_, err = scheduler.AddJob("* * * * *", domria.NewRunner(misc.Secondary, cfg.Mining.Runners.Domria, db))
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("mining: failed to run domria secondary runner, %v", err)
	}
	scheduler.Run()
	if err = db.Close(); err != nil {
		return fmt.Errorf("mining: failed to close the db, %v", err)
	}
	return nil
}
