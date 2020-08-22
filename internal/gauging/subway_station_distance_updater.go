package gauging

import (
	"database/sql"
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging/metrics"
	"time"
)

func NewSubwayStationDistanceUpdater(db *sql.DB) Updater {
	return &subwayStationDistanceUpdater{db}
}

type subwayStationDistanceUpdater struct {
	db *sql.DB
}

func (updater *subwayStationDistanceUpdater) UpdateFlat(flat *dto.Flat, value float64) error {
	start := time.Now()
	_, err := updater.db.Exec(
		"update flats set subway_station_distance = $1 where origin_url = $2",
		value,
		flat.OriginURL,
	)
	metrics.SubwayStationDistanceDuration.WithLabelValues("update").Observe(time.Since(start).Seconds())
	if err != nil {
		metrics.SubwayStationDistanceUpdate.WithLabelValues("failed").Inc()
		return fmt.Errorf("gauging: subway station distance updater failed to update flat, %v", err)
	}
	metrics.SubwayStationDistanceUpdate.WithLabelValues("successful").Inc()
	return nil
}
