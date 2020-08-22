package gauging

import (
	"database/sql"
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging/metrics"
	"time"
)

func NewGreenZoneDistanceUpdater(db *sql.DB) Updater {
	return &greenZoneDistanceUpdater{db}
}

type greenZoneDistanceUpdater struct {
	db *sql.DB
}

func (updater *greenZoneDistanceUpdater) UpdateFlat(flat *dto.Flat, value float64) error {
	start := time.Now()
	_, err := updater.db.Exec(
		"update flats set green_zone_distance = $1 where origin_url = $2",
		value,
		flat.OriginURL,
	)
	metrics.GreenZoneDistanceDuration.WithLabelValues("update").Observe(time.Since(start).Seconds())
	if err != nil {
		metrics.GreenZoneDistanceUpdate.WithLabelValues("failed").Inc()
		return fmt.Errorf("gauging: green zone distance updater failed to update flat, %v", err)
	}
	metrics.GreenZoneDistanceUpdate.WithLabelValues("successful").Inc()
	return nil
}
