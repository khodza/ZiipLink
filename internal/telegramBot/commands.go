package telegrambot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCommand = "start"
	helpCommand  = "help"
)

func (b *Bot) handleCommands(update tgbotapi.Update) error {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case startCommand:
		//set user state to WaitingForURL that means user
		UserStates[update.Message.From.ID] = WaitingForURL
		b.startCommandHandler(&msg)
	case helpCommand:
		b.helpCommandHandler(&msg)
	default:
		msg.Text = "Unknown command"
	}

	//send message
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) startCommandHandler(mes *tgbotapi.MessageConfig) {
	mes.Text = "Welcome to the ZipLink bot! 😅\n\nPlease enter the URL you want to shorten:"
}

func (b *Bot) helpCommandHandler(mes *tgbotapi.MessageConfig) {
	mes.Text = "This is ZipLink bot. To start, enter /start to create a short link."
}
