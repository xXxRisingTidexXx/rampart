package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewMarkupHandler(
	config config.ModeratorHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	return &markupHandler{
		&helper{bot},
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.LuxuryButton),
				tgbotapi.NewKeyboardButton(config.ComfortButton),
				tgbotapi.NewKeyboardButton(config.JunkButton),
				tgbotapi.NewKeyboardButton(config.ConstructionButton),
				tgbotapi.NewKeyboardButton(config.ExcessButton),
			),
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.EnoughButton)),
		),
	}
}

type markupHandler struct {
	helper *helper
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *markupHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "markup"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	var (
		id  int
		url string
	)
	row := tx.QueryRow(
		`select id, url
		from images
		where interior = 'unknown'
			and id not in (select image_id from moderations)
		order by random()
		limit 1`,
	)
	if err := row.Scan(&id, &url); err == sql.ErrNoRows {
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, h.helper.sendMessage(update.Message.Chat.ID, "no_markup", nil)
	} else if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to read an image, %v", err)
	}
	_, err = tx.Exec(
		`insert into moderations (id, image_id) values ($1, $2)`,
		update.Message.Chat.ID,
		id,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to create a moderation, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, h.helper.sendText(update.Message.Chat.ID, url, h.markup)
}
