package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"strconv"
	"strings"
)

func NewLikeHandler(config config.AssistantHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &likeHandler{&helper{bot}, db, config.Separator}
}

type likeHandler struct {
	helper    *helper
	db        *sql.DB
	separator string
}

func (h *likeHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	info := NewInfo("like")
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
	info.Extras["like"] = "check"
	ok, err := h.updateLookup(tx, id, true)
	if err != nil {
		_ = tx.Rollback()
		return info, err
	}
	if ok {
		if err := tx.Commit(); err != nil {
			return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return info, h.helper.answerCallback(update.CallbackQuery.ID, "check_like")
	}
	info.Extras["like"] = "uncheck"
	ok, err = h.updateLookup(tx, id, false)
	if err != nil {
		_ = tx.Rollback()
		return info, err
	}
	if ok {
		if err := tx.Commit(); err != nil {
			return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return info, h.helper.answerCallback(update.CallbackQuery.ID, "uncheck_like")
	}
	info.Extras["like"] = "absent"
	if err := tx.Commit(); err != nil {
		return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return info, h.helper.answerCallback(update.CallbackQuery.ID, "absent_like")
}

func (h *likeHandler) updateLookup(tx *sql.Tx, id int64, shouldLike bool) (bool, error) {
	result, err := tx.Exec(
		`update lookups set is_liked = $1 where id = $2 and is_liked = not $1`,
		shouldLike,
		id,
	)
	if err != nil {
		return false, fmt.Errorf("telegram: handler failed to update a lookup, %v", err)
	}
	number, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("telegram: handler failed to get affected row number, %v", err)
	}
	return number > 0, nil
}
