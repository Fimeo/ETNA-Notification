package etna

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"etna-notification/internal/infrastructure/network"
)

const (
	loginURL       = "https://auth.etna-alternance.net/identity"
	informationURL = "https://intra-api.etna-alternance.net/students/%s/informations"
)

type Service struct {
}

func (s Service) defaultHeaders() map[string]string {
	return map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
		"Content-Type":    "application/json;charset=UTF-8",
	}
}

func (s Service) Login(authentication Authentication) (*http.Cookie, error) {
	bodyJSON, err := json.Marshal(authentication)
	if err != nil {
		return nil, errors.New("[ERROR] Marshal error")
	}

	resp, err := network.Exec(http.MethodPost, loginURL, bodyJSON, s.defaultHeaders(), nil)
	if err != nil {
		return nil, err
	}

	if len(resp.Cookies) == 0 {
		return nil, errors.New("[ERROR] Connection failed, no cookie in response body")
	}

	return resp.Cookies[0], nil
}

func (s Service) RetrieveNotifications(cookie *http.Cookie, username string) ([]Notification, error) {
	resp, err := network.Exec(
		http.MethodGet,
		fmt.Sprintf(informationURL, username), // viper.GetString("etna.user")
		nil,
		nil,
		[]*http.Cookie{cookie})
	if err != nil {
		return nil, err
	}

	var notification []Notification
	err = json.Unmarshal(resp.Body, &notification)
	if err != nil {
		log.Print("[ERROR] UnMarshal error")
		return nil, err
	}

	return notification, nil
}
