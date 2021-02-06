package domria

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewMiner(config config.DomriaMiner, db *sql.DB, logger log.FieldLogger) *Miner {
	return &Miner{}
}

type Miner struct {
}

func (m *Miner) Run() {
	
}
