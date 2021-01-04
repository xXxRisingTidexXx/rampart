package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewDispatcher(
	config config.Dispatcher,
	db *sql.DB,
	logger log.FieldLogger,
) (*Dispatcher, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("telegram: dispatcher failed to instantiate, %v", err)
	}
	return &Dispatcher{bot, config.Timeout}, nil
}

type Dispatcher struct {
	bot     *tgbotapi.BotAPI
	timeout int
}

func (dispatcher *Dispatcher) Pull() {
	updates, _ := dispatcher.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: dispatcher.timeout})
	for update := range updates {

	}
}

//var roomNumberMarkup = tgbotapi.NewInlineKeyboardMarkup(
//	tgbotapi.NewInlineKeyboardRow(
//		tgbotapi.NewInlineKeyboardButtonData("\U0001F937", "any"),
//		tgbotapi.NewInlineKeyboardButtonData("1", "one"),
//		tgbotapi.NewInlineKeyboardButtonData("2", "two"),
//		tgbotapi.NewInlineKeyboardButtonData("3", "three"),
//		tgbotapi.NewInlineKeyboardButtonData("4+", "many"),
//	),
//)

//var floorMarkup = tgbotapi.NewInlineKeyboardMarkup(
//	tgbotapi.NewInlineKeyboardRow(
//		tgbotapi.NewInlineKeyboardButtonData("Байдуже \U0001F612", "any"),
//		tgbotapi.NewInlineKeyboardButtonData("Низький \U0001F60C", "low"),
//		tgbotapi.NewInlineKeyboardButtonData("Високий \U0001F9D0", "high"),
//	),
//)

//updates, _ := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: c.Telegram.Timeout})
//for update := range updates {
//	if update.Message != nil && update.Message.Chat != nil {
//		message := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
//		if update.Message.Command() == "roomnumber" {
//			message.Text = "Скільки кімнат має бути в помешканні? Обирай \U0001F937, якщо це не принципиово."
//			message.ReplyMarkup = roomNumberMarkup
//		} else if update.Message.Command() == "floor" {
//			message.Text = "На якому поверсі хочеш мати домівку?"
//			message.ReplyMarkup = floorMarkup
//		}
//		_, _ = bot.Send(message)
//	}
//	if update.CallbackQuery != nil {
//		entry.Info(update.CallbackQuery.Data)
//		_, _ = bot.DeleteMessage(
//			tgbotapi.NewDeleteMessage(
//				update.CallbackQuery.Message.Chat.ID,
//				update.CallbackQuery.Message.MessageID,
//			),
//		)
//	}
//}
