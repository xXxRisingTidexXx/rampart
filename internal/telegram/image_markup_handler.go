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
	var (
		id  int
		url string
	)
	row := tx.QueryRow(
		`select id, url from images where interior = 'unknown' order by random() limit 1`,
	)
	if err := row.Scan(&id, &url); err == sql.ErrNoRows {
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, h.helper.sendMessage(update.Message.Chat.ID, "no_markup", nil)
	} else if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to select an image, %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, nil
}
