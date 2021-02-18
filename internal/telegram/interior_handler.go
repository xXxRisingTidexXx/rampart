package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewInteriorHandler(config config.ModeratorHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &interiorHandler{
		&helper{bot},
		db,
		map[string]misc.Interior{
			config.LuxuryButton:       misc.LuxuryInterior,
			config.ComfortButton:      misc.ComfortInterior,
			config.JunkButton:         misc.JunkInterior,
			config.ConstructionButton: misc.ConstructionInterior,
			config.ExcessButton:       misc.ExcessInterior,
		},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.MarkupButton)),
		),
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

type interiorHandler struct {
	helper        *helper
	db            *sql.DB
	mappings      map[string]misc.Interior
	absentMarkup  tgbotapi.ReplyKeyboardMarkup
	presentMarkup tgbotapi.ReplyKeyboardMarkup
}

func (h *interiorHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "interior"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	var id int
	row := tx.QueryRow(`select image_id from moderations where id = $1`, update.Message.Chat.ID)
	if err := row.Scan(&id); err == sql.ErrNoRows {
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, nil
	} else if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to read a moderation, %v", err)
	}
	interior, ok := h.mappings[update.Message.Text]
	if !ok {
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, h.helper.sendMessage(update.Message.Chat.ID, "invalid_interior", nil)
	}
	_, err = tx.Exec(`update images set interior = $1 where id = $2`, interior.String(), id)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler faield to update an image, %v", err)
	}
	var url string
	row = tx.QueryRow(
		`select id, url
		from images
		where interior = 'unknown'
			and id not in (select image_id from moderations)
		order by random()
		limit 1`,
	)
	if err := row.Scan(&id, &url); err == sql.ErrNoRows {
		_, err = tx.Exec(`delete from moderations where id = $1`, update.Message.Chat.ID)
		if err != nil {
			_ = tx.Rollback()
			return fields, fmt.Errorf("telegram: handler failed to delete a moderation, %v", err)
		}
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, h.helper.sendMessage(
			update.Message.Chat.ID,
			"absent_interior",
			h.absentMarkup,
		)
	} else if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to read an image, %v", err)
	}
	_, err = tx.Exec(
		`update moderations set image_id = $1 where id = $2`,
		id,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to update a moderation, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, h.helper.sendImage(update.Message.Chat.ID, id, url, h.presentMarkup)
}
