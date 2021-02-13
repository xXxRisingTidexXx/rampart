package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var MessisMinings = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_messis_minings_total",
		Help: "Reflects overall parsed flat quantity",
	},
	[]string{"miner", "status"},
)

var MessisMiningRetries = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_messis_mining_retries_total",
		Help: "Observes data source request attempts",
	},
	[]string{"miner"},
)

var MessisMiningSanitations = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_messis_mining_sanitations_total",
		Help: "Counts various sanitation effect phenomenas",
	},
	[]string{"miner", "action"},
)

var MessisMiningDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_messis_mining_duration_seconds",
		Help:    "Tracks miner workflow timing",
		Buckets: []float64{0.5, 1, 2, 5, 10, 15, 20, 30},
	},
	[]string{"miner"},
)

var MessisGeocodings = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_messis_geocodings_total",
		Help: "Collects position detection cases",
	},
	[]string{"status"},
)

var MessisGeocodingDuration = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "rampart_messis_geocoding_duration_seconds",
		Help:    "Monitors position recognition timing",
		Buckets: []float64{0.5, 1, 2, 5, 10, 15, 20, 30},
	},
)

var MessisGaugings = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_messis_gaugings_total",
		Help: "Tracks geographical feature calculation",
	},
	[]string{"host", "feature", "status"},
)

var MessisGaugingDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_messis_gauging_duration_seconds",
		Help:    "Measures location-based property computation time",
		Buckets: []float64{0.5, 1, 2, 5, 10, 15, 20, 30},
	},
	[]string{"host", "feature"},
)

var MessisStorings = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_messis_storings_total",
		Help: "Holds DB interactions",
	},
	[]string{"resource", "status"},
)

var MessisStoringDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "rampart_messis_storing_duration_seconds",
		Help:    "Observes SQL query durations",
		Buckets: []float64{0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
	},
	[]string{"resource", "action"},
)

var TelegramUpdates = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_telegram_updates_total",
		Help: "Collects incoming Telegram API updates",
	},
	[]string{"handler"},
)
