package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewModeratorHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &moderatorHandler{}
}

type moderatorHandler struct {

}

func (h *moderatorHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if update.Message != nil && update.Message.Chat != nil && update.Message.Chat.UserName == "junkkerrigan" {
		log.Info("Hello, bitch!")
	}
	return nil, nil
}
