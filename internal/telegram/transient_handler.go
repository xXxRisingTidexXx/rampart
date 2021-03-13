package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TransientHandler interface {
	HandleTransientUpdate(tgbotapi.Update, *sql.Tx) (tgbotapi.MessageConfig, error)
}
