package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type StatusHandler interface {
	HandleStatusUpdate(update tgbotapi.Update, tx *sql.Tx) (tgbotapi.MessageConfig, error)
}
