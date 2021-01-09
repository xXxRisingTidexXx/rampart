package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
)

func NewStartHandler(db * sql.DB) Handler {
	return &startHandler{
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Підписка \U0001F49C"),
				tgbotapi.NewKeyboardButton("Довідка \U0001F64B"),
			),
		),
	}
}

type startHandler struct {
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (handler *startHandler) Name() string {
	return "start"
}

// TODO: add message randomization.
func (handler *startHandler) HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	if update.Message == nil ||
		update.Message.Chat == nil ||
		!(update.Message.Command() == handler.Name() ||
			update.Message.Text == "Головне меню \U00002B05") {
		return false, nil
	}
	tx, err := handler.db.Begin()
	if err != nil {
		return true, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/start.html"))
	if err != nil {
		return true, fmt.Errorf("telegram: handler failed to read a file, %v", err)
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = handler.markup
	if _, err := bot.Send(message); err != nil {
		return true, fmt.Errorf("telegram: handler failed to send a message, %v", err)
	}
	return true, nil
}
