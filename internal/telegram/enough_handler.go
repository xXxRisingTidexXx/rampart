package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewEnoughHandler(config config.ModeratorHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &enoughHandler{&helper{bot}, db}
}

type enoughHandler struct {
	helper *helper
	db     *sql.DB
}

func (h *enoughHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "enough"}

	return fields, nil
}
