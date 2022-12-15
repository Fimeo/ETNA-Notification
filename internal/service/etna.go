package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/imroc/req/v3"

	"etna-notification/internal/domain"
)

const (
	loginURL       = "https://auth.etna-alternance.net/identity"
	informationURL = "https://intra-api.etna-alternance.net/students/%s/informations"
)

// etnaWebService is a service that retrieves data from etna web services.
// To retrieve any data, authentication is required. The LoginCookie get the authentication cookie
// to perform other request.
type etnaWebService struct {
	C *req.Client
}

type IEtnaWebService interface {
	LoginCookie(login, password string) (*http.Cookie, error)
	RetrieveUnreadNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error)
	RetrieveAllNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error)
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
	if err != nil {
		return nil, err
	}
	if len(response.Cookies()) == 0 {
		return nil, errors.New("[ERROR] Connection failed, no cookie in response body")
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "authenticator" {
			return cookie, nil
		}
	}

	return nil, errors.New("[ERROR] Connection failed, authenticator cookie not found in response")
}

func (s *etnaWebService) RetrieveUnreadNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error) {
	_, err = s.C.R().
		SetResult(&notifications).
		SetCookies(authenticationCookie).
		Get(fmt.Sprintf(informationURL, username))
	if err != nil {
		return nil, err
	}

	return
}

func (s *etnaWebService) RetrieveAllNotifications(authenticationCookie *http.Cookie, username string) (notifications []*domain.EtnaNotification, err error) {
	_, err = s.C.R().
		SetResult(&notifications).
		SetCookies(authenticationCookie).
		Get(fmt.Sprintf(informationURL, username) + "/archived")
	if err != nil {
		return nil, err
	}

	return
}
