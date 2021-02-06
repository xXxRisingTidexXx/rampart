package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TransientHandler interface {
	HandleTransientUpdate(update tgbotapi.Update, tx *sql.Tx) (tgbotapi.MessageConfig, error)
}
