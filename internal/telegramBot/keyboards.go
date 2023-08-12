package telegrambot

import (
	sl "zipinit/internal/lib/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var randomGenKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Make randomly", "random-alias"),
	),
)

func (b *Bot) handleInlineKeyboard(update tgbotapi.Update) error {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "URL accepted wait a second")
	if _, err := b.bot.Request(callback); err != nil {
		b.log.Error("failed to send callback", sl.Err(err))
		return err
	}

	// And finally, send a message containing the data received.
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	// Handle button callback here
	switch update.CallbackQuery.Data {
	case "random-alias":
		b.handleSceneInput(update, UserStates[update.CallbackQuery.From.ID], b.service)
	}
	//send message
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("failed to send message", sl.Err(err))
		return err
	}
	return nil
}
