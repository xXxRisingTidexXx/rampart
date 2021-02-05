package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewSender(bot *tgbotapi.BotAPI, logger log.FieldLogger) *Sender {
	return &Sender{&helper{bot}, logger}
}

type Sender struct {
	helper *helper
	logger log.FieldLogger
}

func (s *Sender) SendLookup(lookup Lookup) {

}
