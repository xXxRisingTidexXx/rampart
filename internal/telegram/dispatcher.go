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
		[]Handler{
			NewStartHandler(db),
			NewHelpHandler(),
			NewAddHandler(db),
			NewCityHandler(db),
		},
		logger,
	}, nil
}

type Dispatcher struct {
	bot          *tgbotapi.BotAPI
	timeout      int
	workerNumber int
	handlers     []Handler
	logger       log.FieldLogger
}

func (dispatcher *Dispatcher) Pull() {
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
		var (
			ok      = false
			err     error
			handler = "none"
		)
		for i := 0; i < len(dispatcher.handlers) && !ok; i++ {
			ok, err = dispatcher.handlers[i].HandleUpdate(dispatcher.bot, update)
			if ok {
				handler = dispatcher.handlers[i].Name()
			}
			if err != nil {
				dispatcher.logger.WithField("handler", dispatcher.handlers[i].Name()).Error(err)
			}
		}
		metrics.TelegramUpdates.WithLabelValues(handler).Inc()
	}
	group.Done()
}
