package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"sync"
)

func StartDispatcher(
	config config.Dispatcher,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
	logger log.FieldLogger,
) {
	updates, _ := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: config.Timeout})
	handler := NewRootHandler(config.Handler, bot, db)
	group := &sync.WaitGroup{}
	group.Add(config.WorkerNumber)
	for i := 0; i < config.WorkerNumber; i++ {
		go work(updates, handler, logger, group)
	}
	group.Wait()
}

func work(
	updates tgbotapi.UpdatesChannel,
	handler Handler,
	logger log.FieldLogger,
	group *sync.WaitGroup,
) {
	for update := range updates {
		fields, err := handler.HandleUpdate(update)
		if err != nil {
			logger.WithFields(fields).Error(err)
		}
		metrics.TelegramUpdates.WithLabelValues(fields["handler"].(string)).Inc()
	}
	group.Done()
}
