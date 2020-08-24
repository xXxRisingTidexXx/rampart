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

var IndustrialZoneDistanceUpdate = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_industrial_zone_distance_update_count",
		Help: "Counts industrial zone distance update results",
	},
	[]string{"status"},
)

var IndustrialZoneDistanceDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_gauging_industrial_zone_distance_duration_seconds",
		Help:    "Stage duration spent for industrial zone distance tracking",
		Buckets: []float64{0.001, 0.005, 0.05, 0.5, 1, 5, 10, 20, 30},
	},
	[]string{"stage"},
)

var GreenZoneDistanceGauging = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_green_zone_distance_gauging_count",
		Help: "Counts green zone distance gauging results",
	},
	[]string{"status"},
)

var GreenZoneDistanceUpdate = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_gauging_green_zone_distance_update_count",
		Help: "Counts green zone distance update results",
	},
	[]string{"status"},
)

var GreenZoneDistanceDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_gauging_green_zone_distance_duration_seconds",
		Help:    "Stage duration spent for green zone distance tracking",
		Buckets: []float64{0.001, 0.005, 0.05, 0.5, 1, 5, 10, 20, 30},
	},
	[]string{"stage"},
)

func init() {
	SubwayStationDistanceGauging.WithLabelValues("successful")
	SubwayStationDistanceGauging.WithLabelValues("subwayless")
	SubwayStationDistanceGauging.WithLabelValues("inconclusive")
	SubwayStationDistanceGauging.WithLabelValues("failed")
	SubwayStationDistanceUpdate.WithLabelValues("successful")
	SubwayStationDistanceUpdate.WithLabelValues("failed")
	SubwayStationDistanceDuration.WithLabelValues("gauging")
	SubwayStationDistanceDuration.WithLabelValues("update")
	IndustrialZoneDistanceGauging.WithLabelValues("successful")
	IndustrialZoneDistanceGauging.WithLabelValues("subwayless")
	IndustrialZoneDistanceGauging.WithLabelValues("inconclusive")
	IndustrialZoneDistanceGauging.WithLabelValues("failed")
	IndustrialZoneDistanceUpdate.WithLabelValues("successful")
	IndustrialZoneDistanceUpdate.WithLabelValues("failed")
	IndustrialZoneDistanceDuration.WithLabelValues("gauging")
	IndustrialZoneDistanceDuration.WithLabelValues("update")
	GreenZoneDistanceGauging.WithLabelValues("successful")
	GreenZoneDistanceGauging.WithLabelValues("subwayless")
	GreenZoneDistanceGauging.WithLabelValues("inconclusive")
	GreenZoneDistanceGauging.WithLabelValues("failed")
	GreenZoneDistanceUpdate.WithLabelValues("successful")
	GreenZoneDistanceUpdate.WithLabelValues("failed")
	GreenZoneDistanceDuration.WithLabelValues("gauging")
	GreenZoneDistanceDuration.WithLabelValues("update")
}
