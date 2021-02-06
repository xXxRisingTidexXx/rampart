package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewLikeHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &likeHandler{&helper{bot}, db, config.Separator}
}

type likeHandler struct {
	helper    *helper
	db        *sql.DB
	separator string
}

func (h *likeHandler) HandleUpdate(_ tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "like"}, nil
}
