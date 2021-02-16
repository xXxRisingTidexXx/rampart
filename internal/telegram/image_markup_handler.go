package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewImageMarkupHandler(
	config config.ModeratorHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	return &imageMarkupHandler{&helper{bot}, db}
}

type imageMarkupHandler struct {
	helper *helper
	db     *sql.DB
}

func (h *imageMarkupHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "image-markup"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, nil
}
