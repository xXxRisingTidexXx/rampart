package gauging

import (
	"database/sql"
	"github.com/prometheus/common/log"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
)

func NewGreenZoneDistanceUpdater(db *sql.DB) Updater {
	return &greenZoneDistanceUpdater{db}
}

type greenZoneDistanceUpdater struct {
	db *sql.DB
}

func (updater *greenZoneDistanceUpdater) UpdateFlat(flat *dto.Flat, value float64) {
	_, err := updater.db.Exec(
		"update flats set green_zone_distance = $1 where origin_url = $2",
		value,
		flat.OriginURL,
	)
	if err != nil {
		log.Errorf("gauging: updater failed to update flat, %v", err)
	}
}
