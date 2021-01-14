package telegram

//func sendMessage(
//	bot *tgbotapi.BotAPI,
//	update tgbotapi.Update,
//	template string,
//	markup interface{},
//) error {
//	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/" + template + ".html"))
//	if err != nil {
//		return fmt.Errorf("telegram: failed to read a file, %v", err)
//	}
//	message := tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
//	message.ParseMode = tgbotapi.ModeHTML
//	message.ReplyMarkup = markup
//	if _, err := bot.Send(message); err != nil {
//		return fmt.Errorf("telegram: failed to send a message, %v", err)
//	}
//	return nil
//}
