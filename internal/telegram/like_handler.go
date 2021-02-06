package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"strconv"
	"strings"
)

func NewLikeHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &likeHandler{&helper{bot}, db, config.Separator}
}

type likeHandler struct {
	helper    *helper
	db        *sql.DB
	separator string
}

func (h *likeHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "like"}
	_, err := strconv.ParseInt(
		update.CallbackQuery.Data[strings.LastIndex(update.CallbackQuery.Data, h.separator)+1:],
		10,
		64,
	)
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to parse id, %v", err)
	}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, nil
}
