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
		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 15, 20, 30},
	},
	[]string{"miner"},
)

var TelegramUpdates = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_telegram_updates_total",
		Help: "Collects incoming Telegram API updates",
	},
	[]string{"handler"},
)
