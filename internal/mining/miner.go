package mining

import (
	"github.com/robfig/cron/v3"
)

type Miner interface {
	cron.Job
	Spec() string
	Port() int
}
