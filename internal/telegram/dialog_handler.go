package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewDialogHandler(config config.AssistantHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &dialogHandler{
		bot,
		db,
		map[misc.Status]TransientHandler{
			misc.CityStatus:       NewCityHandler(config, bot, db),
			misc.PriceStatus:      NewPriceHandler(config, bot, db),
			misc.RoomNumberStatus: NewRoomNumberHandler(config, bot, db),
			misc.FloorStatus:      NewFloorHandler(config, bot, db),
		},
	}
}

type dialogHandler struct {
	bot      *tgbotapi.BotAPI
	db       *sql.DB
	handlers map[misc.Status]TransientHandler
}

func (h *dialogHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	info := NewInfo("dialog")
	tx, err := h.db.Begin()
	if err != nil {
		return info, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	var view string
	row := tx.QueryRow(`select status from transients where id = $1`, update.Message.Chat.ID)
	if err := row.Scan(&view); err == sql.ErrNoRows {
		if err := tx.Commit(); err != nil {
			return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return info, nil
	} else if err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to read a transient, %v", err)
	}
	info.Extras["status"] = view
	status, ok := misc.ToStatus(view)
	if !ok {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to find a status")
	}
	handler, ok := h.handlers[status]
	if !ok {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to handle a status")
	}
	info = NewInfo(view)
	message, err := handler.HandleTransientUpdate(update, tx)
	if err != nil {
		_ = tx.Rollback()
		return info, err
	}
	if err := tx.Commit(); err != nil {
		return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	if _, err := h.bot.Send(message); err != nil {
		return info, fmt.Errorf("telegram: handler failed to send a message, %v", err)
	}
	return info, nil
}
