package telegrambot

import (
	"log"
	sl "zipinit/internal/lib/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slog"
)

type Service interface {
	SaveUrl(url, alias string) (string, error)
	GetUrl(alias string) (string, error)
}

type Bot struct {
	bot     *tgbotapi.BotAPI
	service Service
	log     *slog.Logger
}

func NewBot(token string, service Service, log *slog.Logger) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error("failed to init bot", sl.Err(err))
	}

	bot.Debug = true
	return &Bot{
		bot:     bot,
		service: service,
		log:     log,
	}
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdates()

	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		var userID int64
		if update.Message != nil {
			userID = update.Message.From.ID
			state := UserStates[userID]

			if update.Message.IsCommand() {
				b.handleCommands(update)
				continue
			}

			b.handleSceneInput(update, state, b.service)

		} else if update.CallbackQuery != nil {
			userID = update.CallbackQuery.From.ID
			state := UserStates[userID]
			_ = state

			b.handleInlineKeyboard(update)
		}

	}
}

func (b *Bot) initUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
