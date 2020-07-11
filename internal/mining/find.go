package mining

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"rampart/internal/mining/domria"
	"rampart/internal/mining/metrics"
)

func FindMiner(
	alias string,
	config *config.Miners,
	db *sql.DB,
	gatherer *metrics.Gatherer,
	logger log.FieldLogger,
) (Miner, error) {
	miners := map[string]Miner{
		config.DomriaPrimary.Alias:   domria.NewMiner(config.DomriaPrimary, db, gatherer, logger),
		config.DomriaSecondary.Alias: domria.NewMiner(config.DomriaSecondary, db, gatherer, logger),
	}
	if miner, ok := miners[alias]; ok {
		return miner, nil
	}
	return nil, fmt.Errorf("mining: failed to find the miner with the alias %s", alias)
}
