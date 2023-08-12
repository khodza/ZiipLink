package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"zipinit/internal/config"
	"zipinit/internal/http-server/handlers/redirect"
	"zipinit/internal/http-server/handlers/save"
	"zipinit/internal/http-server/middleware/auth"
	"zipinit/internal/lib/logger"
	sl "zipinit/internal/lib/logger"

	"zipinit/internal/services"
	mongodb "zipinit/internal/storage/mongoDB"
	telegrambot "zipinit/internal/telegramBot"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	//init database
	dataBaseName := cfg.DataBaseName
	storage, err := mongodb.NewStorage(cfg.MongoDBUrl, dataBaseName)

	//init service
	service := services.NewService(storage, log, cfg)

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	//init router
	router := mux.NewRouter()

	saveUrlHandler := save.New(log, service)
	redirectHandler := redirect.New(log, service)

	securedURLHandler := auth.BasicAuthMiddleware(saveUrlHandler, cfg.User, cfg.Password)

	router.Handle("/url", securedURLHandler).Methods(http.MethodPost)

	router.HandleFunc("/{alias}", redirectHandler).Methods(http.MethodGet)

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.TimeOut,
		WriteTimeout: cfg.TimeOut,
		IdleTimeout:  cfg.IdleTimeOut,
	}

	//Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Info("shutting down server gracefully...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Error("error while shutting down server", sl.Err(err))
		}
	}()

	//init telegram bot
	bot := telegrambot.NewBot(cfg.TelegramBotToken, service, log)
	go bot.Start()

	//start server
	err = srv.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", sl.Err(err))
	}
	log.Error("server stopped")
}
