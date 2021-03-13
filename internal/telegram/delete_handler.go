package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"strconv"
	"strings"
)

func NewDeleteHandler(config config.AssistantHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &deleteHandler{&helper{bot}, db, config.Separator}
}

type deleteHandler struct {
	helper    *helper
	db        *sql.DB
	separator string
}

func (h *deleteHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	info := NewInfo("delete")
	id, err := strconv.ParseInt(
		update.CallbackQuery.Data[strings.LastIndex(update.CallbackQuery.Data, h.separator)+1:],
		10,
		64,
	)
	if err != nil {
		return info, fmt.Errorf("telegram: handler failed to parse id, %v", err)
	}
	info.Extras["id"] = id
	tx, err := h.db.Begin()
	if err != nil {
		return info, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	result, err := tx.Exec(
		`update subscriptions set status = 'inactive' where id = $1 and status = 'active'`,
		id,
	)
	if err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to update a subscription, %v", err)
	}
	number, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to get affected row number, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	if number == 0 {
		return info, h.helper.answerCallback(update.CallbackQuery.ID, "absent_delete")
	}
	return info, h.helper.answerCallback(update.CallbackQuery.ID, "present_delete")
}
