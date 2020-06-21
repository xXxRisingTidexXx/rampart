package mining

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func Schedule() error {
	scheduler := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	_, err := scheduler.AddFunc(
		"* * * * *",
		func() {
			log.Debugf("mining: hello, fucking world, %s", time.Now())
			panic("mining: run away")
		},
	)
	if err != nil {
		return err
	}
	scheduler.Run()
	return nil
}
