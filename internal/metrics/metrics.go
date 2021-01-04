package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TelegramCommands = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rampart_telegram_commands_total",
		Help: "Collects incoming Telegram command orders",
	},
	[]string{"command"},
)
