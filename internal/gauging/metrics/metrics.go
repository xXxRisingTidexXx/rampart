package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var SubwayStationDistanceGauging = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_subway_station_distance_gauging_count",
		Help: "Counts subway station distance gauging results",
	},
	[]string{"status"},
)

var SubwayStationDistanceUpdate = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_subway_station_distance_update_count",
		Help: "Counts subway station distance update results",
	},
	[]string{"status"},
)

var SubwayStationDistanceDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_gauging_subway_station_distance_duration_seconds",
		Help:    "Stage duration spent for subway station distance tracking",
		Buckets: []float64{0.001, 0.005, 0.05, 0.5, 1, 5, 10, 20, 30},
	},
	[]string{"stage"},
)

var IndustrialZoneDistanceGauging = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_industrial_zone_distance_gauging_count",
		Help: "Counts industrial zone distance gauging results",
	},
	[]string{"status"},
)

var GreenZoneDistanceGauging = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_green_zone_distance_gauging_count",
		Help: "Counts green zone distance gauging results",
	},
	[]string{"status"},
)

//var UpdaterQueryDuration = promauto.NewHistogramVec(
//	prometheus.HistogramOpts{
//		Name:    "rampart_gauging_updater_query_duration_seconds",
//		Help:    "Database SQL query timing",
//		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
//	},
//	[]string{"feature"},
//)
