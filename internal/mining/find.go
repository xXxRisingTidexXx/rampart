package mining

import (
	"database/sql"
	"fmt"
	"rampart/internal/config"
	"rampart/internal/mining/domria"
	"rampart/internal/mining/metrics"
)

func FindMiner(
	alias string,
	config *config.Miners,
	db *sql.DB,
	gatherer *metrics.Gatherer,
) (Miner, error) {
	miners := []Miner{
		domria.NewMiner(config.DomriaPrimary, db, gatherer),
		domria.NewMiner(config.DomriaSecondary, db, gatherer),
	}
	for _, miner := range miners {
		if miner.Alias() == alias {
			return miner, nil
		}
	}
	return nil, fmt.Errorf("mining: failed to find the miner with the alias %s", alias)
}
