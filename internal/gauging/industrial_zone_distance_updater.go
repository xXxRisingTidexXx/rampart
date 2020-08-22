package gauging

import (
	"database/sql"
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging/metrics"
	"time"
)

func NewIndustrialZoneDistanceUpdater(db *sql.DB) Updater {
	return &industrialZoneDistanceUpdater{db}
}

type industrialZoneDistanceUpdater struct {
	db *sql.DB
}

func (updater *industrialZoneDistanceUpdater) UpdateFlat(flat *dto.Flat, value float64) error {
	start := time.Now()
	_, err := updater.db.Exec(
		"update flats set industrial_zone_distance = $1 where origin_url = $2",
		value,
		flat.OriginURL,
	)
	metrics.IndustrialZoneDistanceDuration.WithLabelValues("update").Observe(time.Since(start).Seconds())
	if err != nil {
		metrics.IndustrialZoneDistanceUpdate.WithLabelValues("failed").Inc()
		return fmt.Errorf("gauging: industrial zone distance updater failed to update flat, %v", err)
	}
	metrics.IndustrialZoneDistanceUpdate.WithLabelValues("successful").Inc()
	return nil
}
