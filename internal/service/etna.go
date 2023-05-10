package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/imroc/req/v3"

	"etna-notification/internal/domain"
)

const (
	loginURL       = "https://auth.etna-alternance.net/identity"
	informationURL = "https://intra-api.etna-alternance.net/students/%s/informations"
	eventURL       = "https://intra-api.etna-alternance.net/students/%s/events?start=%s&end=%s"
)

// etnaWebService is a service that retrieves data from etna web services.
// To retrieve any data, authentication is required. The LoginCookie get the authentication cookie
// to perform other request.
type etnaWebService struct {
	C *req.Client
}

type ErrorWrongCredentials error

type IEtnaWebService interface {
	LoginCookie(login, password string) (*http.Cookie, error)
	RetrieveUnreadNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error)
	RetrieveAllNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error)
	RetrieveCalendarEventInRange(authenticationCookie *http.Cookie, username string, start, end time.Time) (events []*domain.EtnaCalendarEvent, err error)
}

func NewEtnaWebservice(client *req.Client) IEtnaWebService {
	return &etnaWebService{
		C: client,
	}
}

func (s *etnaWebService) LoginCookie(login, password string) (*http.Cookie, error) {
	type body struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	response, err := s.C.R().
		SetBody(body{
			Login:    login,
			Password: password,
		}).
		Post(loginURL)
	log.Printf("[INFO] Login response time : %s", response.TotalTime().String())
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusUnauthorized {
		return nil, ErrorWrongCredentials(fmt.Errorf("[ERROR] Wrong credentials for user : %s", login))
	}
	if len(response.Cookies()) == 0 {
		return nil, fmt.Errorf("[ERROR] Connection failed, no cookie in response body, user : %s", login)
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "authenticator" {
			return cookie, nil
		}
	}

	return nil, fmt.Errorf("[ERROR] Connection failed, authenticator cookie not found in response, user : %s", login)
}

func (s *etnaWebService) RetrieveUnreadNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error) {
	response, err := s.C.R().
		SetSuccessResult(&notifications).
		SetCookies(authenticationCookie).
		Get(fmt.Sprintf(informationURL, username))
	log.Printf("[INFO] Retrieve unread notifications response time : %s", response.TotalTime().String())
	return
}

func (s *etnaWebService) RetrieveAllNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error) {
	response, err := s.C.R().
		SetSuccessResult(&notifications).
		SetCookies(authenticationCookie).
		Get(fmt.Sprintf(informationURL, username) + "/archived")
	log.Printf("[INFO] Retrieve all notifications response time : %s", response.TotalTime().String())
	return
}

func (s *etnaWebService) RetrieveCalendarEventInRange(authenticationCookie *http.Cookie, username string, start, end time.Time) (events []*domain.EtnaCalendarEvent, err error) {
	response, err := s.C.R().
		SetSuccessResult(&events).
		SetCookies(authenticationCookie).
		Get(fmt.Sprintf(eventURL, username, start.Format("2006-01-02"), end.Format("2006-01-02")))
	log.Printf("[INFO] Retrieve events response time : %s", response.TotalTime().String())
	return
}
