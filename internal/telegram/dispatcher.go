package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"sync"
)

func NewDispatcher(
	config config.Dispatcher,
	db *sql.DB,
	logger log.FieldLogger,
) (*Dispatcher, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("telegram: dispatcher failed to instantiate, %v", err)
	}
	return &Dispatcher{
		bot,
		config.Timeout,
		config.WorkerNumber,
		NewRootHandler(config.Handler, bot, db),
		logger,
	}, nil
}

type Dispatcher struct {
	bot          *tgbotapi.BotAPI
	timeout      int
	workerNumber int
	handler      Handler
	logger       log.FieldLogger
}

func (dispatcher *Dispatcher) Dispatch() {
	updates, _ := dispatcher.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: dispatcher.timeout})
	group := &sync.WaitGroup{}
	group.Add(dispatcher.workerNumber)
	for i := 0; i < dispatcher.workerNumber; i++ {
		go dispatcher.work(updates, group)
	}
	group.Wait()
}

func (dispatcher *Dispatcher) work(updates tgbotapi.UpdatesChannel, group *sync.WaitGroup) {
	for update := range updates {
		fields, err := dispatcher.handler.HandleUpdate(update)
		if err != nil {
			dispatcher.logger.WithFields(fields).Error(err)
		}
		metrics.TelegramUpdates.WithLabelValues(fields["handler"].(string)).Inc()
	}
	group.Done()
}
