package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewAddHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &addHandler{
		&helper{bot},
		db,
		"Київ",
	}
}

type addHandler struct {
	helper      *helper
	db          *sql.DB
	defaultCity string
}

// TODO: move default city to config.
// TODO: move default price to config.
func (h *addHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "add"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
	}
	_, err = tx.Exec(
		`insert into transients
		(id, status, city, price, room_number, floor)
		values
		($1, $2, $3, 0, 'any', 'any')`,
		update.Message.Chat.ID,
		misc.CityStatus.String(),
		h.defaultCity,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to create a transient, %v", err)
	}
	buttons, city := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(h.defaultCity)), ""
	row := tx.QueryRow(
		`select city from flats where city != $1 group by city order by count(*) desc limit 1`,
		h.defaultCity,
	)
	switch err := row.Scan(&city); err {
	case nil:
		buttons = append(buttons, tgbotapi.NewKeyboardButton(city))
	case sql.ErrNoRows:
	default:
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to select a city, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, h.helper.sendMessage(
		update,
		"add",
		tgbotapi.NewReplyKeyboard(
			buttons,
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	)
}
