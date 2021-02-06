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

func NewDeleteHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &deleteHandler{&helper{bot}, db, config.Separator}
}

type deleteHandler struct {
	helper    *helper
	db        *sql.DB
	separator string
}

func (h *deleteHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "delete"}
	id, err := strconv.ParseInt(
		update.CallbackQuery.Data[strings.LastIndex(update.CallbackQuery.Data, h.separator)+1:],
		10,
		64,
	)
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to parse id, %v", err)
	}
	fields["subscription_id"] = id
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	result, err := tx.Exec(
		`update subscriptions set status = 'inactive' where id = $1 and status = 'active'`,
		id,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to update a subscription, %v", err)
	}
	number, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to get affected row number, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	if number == 0 {
		return fields, h.helper.answerCallback(update.CallbackQuery.ID, "absent_delete")
	}
	return fields, h.helper.answerCallback(update.CallbackQuery.ID, "present_delete")
}
