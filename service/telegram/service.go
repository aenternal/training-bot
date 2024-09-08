package telegram

import (
	tele "gopkg.in/telebot.v3"
)

type AuthService interface {
	GetGoogleOAuthURL(userId int64) string
	GetYandexOAuthURL() string
}

type TelegramService struct {
	bot         *tele.Bot
	authService AuthService
}

func NewTelegramService(b *tele.Bot, service AuthService) *TelegramService {
	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/authYandex", func(c tele.Context) error {
		return c.Send(service.GetYandexOAuthURL())
	})

	b.Handle("/authGoogle", func(c tele.Context) error {
		return c.Send(service.GetGoogleOAuthURL(c.Sender().ID))
	})

	b.Handle("/donate", func(c tele.Context) error {
		inv := tele.Invoice{
			Title:       "Donate",
			Description: "Donate",
			Payload:     "{}",
			Currency:    "XTR",
			Prices: []tele.Price{
				tele.Price{
					Label:  "Donate",
					Amount: 10,
				},
			},
		}

		_, err := inv.Send(b, &tele.User{ID: c.Sender().ID}, nil)
		if err != nil {
			return err
		}

		return nil
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Text received")
	})

	return &TelegramService{
		bot:         b,
		authService: service,
	}
}

func (s *TelegramService) AddHandler(handleText, handler func(c tele.Context) error) {
	s.bot.Handle(handleText, handler)
}

func (s *TelegramService) Start() {
	s.bot.Start()
}

func (s *TelegramService) Stop() {
	s.bot.Stop()
}
