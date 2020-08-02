package mining

import (
	"database/sql"
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/domria"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
)

func FindMiner(
	alias string,
	config *config.Miners,
	db *sql.DB,
	gatherer *metrics.Gatherer,
	logger *logging.Logger,
) (Miner, error) {
	miners := map[string]Miner{
		config.DomriaPrimary.Alias:   domria.NewMiner(config.DomriaPrimary, db, gatherer, logger),
		config.DomriaSecondary.Alias: domria.NewMiner(config.DomriaSecondary, db, gatherer, logger),
	}
	if miner, ok := miners[alias]; ok {
		return miner, nil
	}
	return nil, fmt.Errorf("mining: failed to find the miner")
}
