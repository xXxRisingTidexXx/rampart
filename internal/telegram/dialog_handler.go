package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewDialogHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &dialogHandler{
		bot,
		db,
		map[misc.Status]StatusHandler{
			misc.CityStatus:       NewCityStatusHandler(bot, db),
			misc.PriceStatus:      NewPriceStatusHandler(bot, db),
			misc.RoomNumberStatus: NewRoomNumberStatusHandler(bot, db),
		},
	}
}

type dialogHandler struct {
	bot      *tgbotapi.BotAPI
	db       *sql.DB
	handlers map[misc.Status]StatusHandler
}

func (h *dialogHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "dialog"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	var view string
	row := tx.QueryRow(`select status from transients where id = $1`, update.Message.Chat.ID)
	if err := row.Scan(&view); err == sql.ErrNoRows {
		if err := tx.Commit(); err != nil {
			return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return fields, nil
	} else if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to read a transient, %v", err)
	}
	fields["status"] = view
	status, ok := misc.ToStatus(view)
	if !ok {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to find a status")
	}
	handler, ok := h.handlers[status]
	if !ok {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to handle a status")
	}
	fields = log.Fields{"handler": view}
	message, err := handler.HandleStatusUpdate(update, tx)
	if err != nil {
		_ = tx.Rollback()
		return fields, err
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	if _, err := h.bot.Send(message); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to send a message, %v", err)
	}
	return fields, nil
}
