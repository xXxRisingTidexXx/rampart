package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"strconv"
)

func NewSender(config config.Sender, bot *tgbotapi.BotAPI, logger log.FieldLogger) *Sender {
	return &Sender{&helper{bot}, config.LikeButton, config.LikeAction, config.Separator, logger}
}

type Sender struct {
	helper    *helper
	button    string
	action    string
	separator string
	logger    log.FieldLogger
}

func (s *Sender) SendLookup(lookup Lookup) {
	err := s.helper.sendTemplate(
		lookup.ChatID,
		"lookup",
		lookup,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					s.button,
					s.action+s.separator+strconv.Itoa(lookup.ID),
				),
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
