package telegrambot

import (
	"zipinit/internal/http-server/handlers/save"
	resp "zipinit/internal/lib/api/response"
	sl "zipinit/internal/lib/logger"

	"github.com/go-playground/validator/v10"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState int

const (
	Idle UserState = iota
	WaitingForURL
	WaitingForAlias
	Processing
)

var UserStates = map[int64]UserState{}
var req save.Request

func (b *Bot) handleSceneInput(update tgbotapi.Update, state UserState, service Service) {
	var msg tgbotapi.MessageConfig
	var userID int64

	if update.Message != nil {
		userID = update.Message.From.ID
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.From.ID
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	}

	switch state {
	case WaitingForURL:
		req.URL = update.Message.Text
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			b.log.Error("invalid request", sl.Err(err))
			res := resp.ValidationError(validateErr)
			msg.Text = res.Error
			break
		}
		msg.Text = "Enter alias for the URL:"
		msg.ReplyMarkup = randomGenKeyboard
		UserStates[userID] = WaitingForAlias

	case WaitingForAlias:
		if update.Message != nil {
			req.Alias = update.Message.Text
		}
		//otherwise left req.Alias empty

		UserStates[userID] = Processing

		shortLink, err := service.SaveUrl(req.URL, req.Alias)
		if err != nil {
			b.log.Error("failed to save url", sl.Err(err))
			msg.Text = "Failed to save url"
			break
		}
		msg.Text = "Here is the your short link:"
		//send mess
		_, err = b.bot.Send(msg)
		if err != nil {
			b.log.Error("failed to send message", sl.Err(err))
			break
		}
		msg.Text = shortLink
		delete(UserStates, userID)
	default:
		msg.Text = "Unknown input"
	}

	//send message
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("failed to send message", sl.Err(err))
		return
	}

}
