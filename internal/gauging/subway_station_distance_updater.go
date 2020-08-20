package gauging

import (
	"database/sql"
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
)

func NewSubwayStationDistanceUpdater(db *sql.DB) Updater {
	return &subwayStationDistanceUpdater{db}
}

type subwayStationDistanceUpdater struct {
	db *sql.DB
}

func (updater *subwayStationDistanceUpdater) UpdateFlat(flat *dto.Flat, value float64) error {
	_, err := updater.db.Exec(
		"update flats set subway_station_distance = $1 where origin_url = $2",
		value,
		flat.OriginURL,
	)
	if err != nil {
		return fmt.Errorf("gauging: updater failed to update flat's subway station distance, %v", err)
	}
	return nil
}
