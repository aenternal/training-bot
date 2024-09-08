package auth

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
	"google.golang.org/api/calendar/v3"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"strconv"
)

var (
	oauth2ConfigYandex = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Endpoint:     yandex.Endpoint,
		RedirectURL:  "http://localhost:8080/yandexoauth2callback",
		Scopes:       []string{"calendar:all"},
	}
)

type AuthService struct {
	bot *tele.Bot
}

func NewAuthService(b *tele.Bot) *AuthService {
	return &AuthService{
		bot: b,
	}
}

func (a *AuthService) getOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("http://localhost:8080/oauth2callback"),
		Scopes:       []string{calendar.CalendarReadonlyScope},
	}
}
func (a *AuthService) GetGoogleOAuthURL(userId int64) string {
	return a.getOauthConfig().AuthCodeURL(fmt.Sprintf("%d", userId))
}

func (a *AuthService) GetYandexOAuthURL() string {
	return oauth2ConfigYandex.AuthCodeURL("state")
}

func (a *AuthService) HandleAuthGoogleCallback(e echo.Context) error {
	userId := e.FormValue("state")
	userIdInt64, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}

	code := e.FormValue("code")
	token, err := a.getOauthConfig().Exchange(e.Request().Context(), code)
	if err != nil {
		return e.NoContent(http.StatusInternalServerError)
	}
	log.Info(token.AccessToken)

	_, err = a.bot.Send(&tele.User{ID: userIdInt64}, "spasibo")
	if err != nil {
		return e.NoContent(http.StatusInternalServerError)
	}

	//
	//client := oauth2Config.Client(e.Request().Context(), token)
	//srv, err := calendar.NewService(e.Request().Context(), option.WithHTTPClient(client))
	//if err != nil {
	//	return e.NoContent(http.StatusInternalServerError)
	//}

	//calendarID := "primary"
	//events, err := srv.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).
	//	TimeMin(time.Now().Format(time.RFC3339)).OrderBy("startTime").Do()
	//if err != nil {
	//	return e.NoContent(http.StatusInternalServerError)
	//}
	//
	//log.Infof("Got %d events", len(events.Items))
	//log.Infof("Got %+v events", events.Items)
	//
	//for _, event := range events.Items {
	//	log.Info(event.Summary)
	//	log.Infof(event.Description)
	//	log.Infof(event.Id)
	//	log.Infof(event.Start.DateTime)
	//}
	return e.NoContent(http.StatusOK)
}

func (a *AuthService) HandleAuthYandex(e echo.Context) error {
	url := oauth2ConfigYandex.AuthCodeURL("state")
	return e.Redirect(http.StatusTemporaryRedirect, url)
}
