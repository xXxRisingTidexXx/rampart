package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewPublisher(
	config config.Publisher,
	db *sql.DB,
	logger log.FieldLogger,
) (*Publisher, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("telegram: publisher failed to instantiate, %v", err)
	}
	return &Publisher{NewObserver(db, logger), NewSender(bot, logger)}, nil
}

type Publisher struct {
	observer *Observer
	sender   *Sender
}

func (p *Publisher) Run() {
	for _, lookup := range p.observer.ObserveLookups() {
		p.sender.SendLookup(lookup)
	}
}
