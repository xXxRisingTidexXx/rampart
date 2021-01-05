package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TelegramUpdates = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_telegram_updates_total",
		Help: "Collects incoming Telegram API updates",
	},
	[]string{"handler"},
)
