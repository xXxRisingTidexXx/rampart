package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewListHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &listHandler{
		&helper{bot},
		db,
		config.AnyPricePlaceholder,
		map[string]string{
			misc.AnyRoomNumber.String():   config.AnyRoomNumberPlaceholder,
			misc.OneRoomNumber.String():   config.OneRoomNumberPlaceholder,
			misc.TwoRoomNumber.String():   config.TwoRoomNumberPlaceholder,
			misc.ThreeRoomNumber.String(): config.ThreeRoomNumberPlaceholder,
			misc.ManyRoomNumber.String():  config.ManyRoomNumberPlaceholder,
		},
		map[string]string{
			misc.AnyFloor.String():  config.AnyFloorPlaceholder,
			misc.LowFloor.String():  config.LowFloorPlaceholder,
			misc.HighFloor.String(): config.HighFloorPlaceholder,
		},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.StartButton)),
		),
	}
}

type listHandler struct {
	helper                 *helper
	db                     *sql.DB
	anyPricePlaceholder    string
	roomNumberPlaceholders map[string]string
	floorPlaceholders      map[string]string
	markup                 tgbotapi.ReplyKeyboardMarkup
}

func (h *listHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "list"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	rows, err := tx.Query(
		`select id, uuid, city, price, room_number, floor
		from subscriptions
		where status = 'actual'
			and chat_id = $1`,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to read subscriptions, %v", err)
	}
	subscriptions := make([]subscription, 0)
	for rows.Next() {
		var (
			id         int
			uuid       string
			city       string
			price      float32
			roomNumber string
			floor      string
		)
		err := rows.Scan(&id, &uuid, &city, &price, &roomNumber, &floor)
		if err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			fields["subscription_id"] = id
			return fields, fmt.Errorf("telegram: handler failed to scan a row, %v", err)
		}
		shape := h.anyPricePlaceholder
		if price > 0 {
			shape = fmt.Sprintf("%.2f $", price)
		}
		roomNumber, ok := h.roomNumberPlaceholders[roomNumber]
		if !ok {
			roomNumber = h.roomNumberPlaceholders[misc.AnyRoomNumber.String()]
		}
		floor, ok = h.floorPlaceholders[floor]
		if !ok {
			floor = h.floorPlaceholders[misc.AnyFloor.String()]
		}
		subscriptions = append(
			subscriptions,
			subscription{id, uuid, city, shape, roomNumber, floor},
		)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to finish iteration, %v", err)
	}
	if err := rows.Close(); err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to close rows, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	if len(subscriptions) == 0 {
		return fields, h.helper.sendMessage(update, "empty_list", h.markup)
	}
	for _, s := range subscriptions {
		if err := h.helper.sendTemplate(update, "full_list", s, h.markup); err != nil {
			fields["subscription_id"] = s.ID
			return fields, err
		}
	}
	return fields, nil
}
