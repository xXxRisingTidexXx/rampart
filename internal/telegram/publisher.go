package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewPublisher(
	config config.Publisher,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
	logger log.FieldLogger,
) *Publisher {
	return &Publisher{
		NewObserver(db, logger),
		NewSender(config.Sender, bot, logger),
		NewReviewer(db, logger),
	}
}

type Publisher struct {
	observer *Observer
	sender   *Sender
	reviewer *Reviewer
}

func (p *Publisher) Run() {
	for _, lookup := range p.observer.ObserveLookups() {
		p.sender.SendLookup(lookup)
		p.reviewer.ReviewLookup(lookup)
	}
}
