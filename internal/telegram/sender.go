package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func NewSender(bot *tgbotapi.BotAPI, logger log.FieldLogger) *Sender {
	return &Sender{&helper{bot}, logger}
}

type Sender struct {
	helper *helper
	logger log.FieldLogger
}

func (s *Sender) SendLookup(lookup Lookup) {
	err := s.helper.sendTemplate(
		lookup.ChatID,
		"lookup",
		lookup,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Охуєнно", "like|"+strconv.Itoa(lookup.ID)),
			),
		),
	)
	if err != nil {
		s.logger.WithField("lookup_id", lookup.ID).Errorf(
			"telegram: sender failed to send a lookup, %v",
			err,
		)
	}
}
