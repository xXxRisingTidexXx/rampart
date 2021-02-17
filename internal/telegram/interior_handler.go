package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewInteriorHandler(config config.ModeratorHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &interiorHandler{&helper{bot}, db}
}

type interiorHandler struct {
	helper *helper
	db     *sql.DB
}

func (h *interiorHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "interior"}

	return fields, nil
}
