package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewListHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &listHandler{&helper{bot}, db}
}

type listHandler struct {
	helper *helper
	db     *sql.DB
}

func (h *listHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "list"}

	return fields, nil
}
