package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var SubwayStationDistance = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_subway_station_distance_count",
		Help: "Counts subway station distance gauging results",
	},
	[]string{"status"},
)

var IndustrialZoneDistance = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_industrial_zone_distance_count",
		Help: "Counts industrial zone distance gauging results",
	},
	[]string{"status"},
)

var GreenZoneDistance = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_green_zone_distance_count",
		Help: "Counts green zone distance gauging results",
	},
	[]string{"status"},
)

var OverpassRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_gauging_overpass_request_duration_seconds",
		Help:    "Overpass API HTTP request timing",
		Buckets: []float64{5, 10, 15, 20, 30},
	},
	[]string{"feature"},
)

var DBQueryDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_gauging_db_query_duration_seconds",
		Help:    "Database SQL query timing",
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
	},
	[]string{"feature"},
)
