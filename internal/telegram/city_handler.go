package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewCityHandler(db *sql.DB) Handler {
	return &cityHandler{db}
}

type cityHandler struct {
	db *sql.DB
}

func (handler *cityHandler) Name() string {
	return "city"
}

func (handler *cityHandler) ServeUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	return false, nil
}
