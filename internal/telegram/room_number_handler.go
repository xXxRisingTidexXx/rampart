package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewRoomNumberHandler(db *sql.DB) Handler {
	return &roomNumberHandler{db}
}

type roomNumberHandler struct {
	db *sql.DB
}

func (handler *roomNumberHandler) Name() string {
	return "room-number"
}

func (handler *roomNumberHandler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) (bool, error) {
	return true, nil
}
