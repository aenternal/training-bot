package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	tele "gopkg.in/telebot.v3"
	"time"
	"training-bot/service/auth"

	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"training-bot/service/telegram"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("starting webinarTelegramBot")
	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	authService := auth.NewAuthService(b)
	telegramService := telegram.NewTelegramService(b, authService)

	go func() {
		telegramService.Start()
	}()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.GET("/authGoogle", authService.HandleAuthGoogle)
	e.GET("/authYandex", authService.HandleAuthYandex)
	e.GET("/oauth2callback", authService.HandleAuthGoogleCallback)

	go func() {
		e.Logger.Fatal(e.Start(":8080"))
	}()

	<-done

	slog.Info("stopping webinarTelegramBot")
	telegramService.Stop()
	e.Close()
}
