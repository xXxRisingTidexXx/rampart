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

// TODO: add templates.
func (h *likeHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "like"}
	id, err := strconv.ParseInt(
		update.CallbackQuery.Data[strings.LastIndex(update.CallbackQuery.Data, h.separator)+1:],
		10,
		64,
	)
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to parse id, %v", err)
	}
	fields["id"] = id
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	fields["reaction"] = "positive"
	ok, err := h.updateLookup(tx, id, true)
	if err != nil {
		_ = tx.Rollback()
		return fields, err
	}
	if ok {
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, h.helper.answerCallback(update.CallbackQuery.ID, "positive_reaction")
	}
	fields["reaction"] = "negative"
	ok, err = h.updateLookup(tx, id, false)
	if err != nil {
		_ = tx.Rollback()
		return fields, err
	}
	if ok {
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, h.helper.answerCallback(update.CallbackQuery.ID, "negative_reaction")
	}
	fields["reaction"] = "absent"
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, h.helper.answerCallback(update.CallbackQuery.ID, "absent_reaction")
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
