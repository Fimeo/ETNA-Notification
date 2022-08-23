package network

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Cookies []*http.Cookie
	Body    []byte
}

// Exec execute HTTP requests according to following parameters.
func Exec(method, url string, body []byte, headers map[string]string, cookies []*http.Cookie) (*Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Print("[ERROR] Request error", err)
		return nil, err
	}

	for key, header := range headers {
		req.Header.Add(key, header)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Print("[ERROR] Request error", err)
		return nil, err
	}
	defer res.Body.Close()

	bodyJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("[ERROR] Request error", err)
		fmt.Println(res.StatusCode, string(bodyJSON))
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(string(bodyJSON))
	}

	return &Response{
		Cookies: res.Cookies(),
		Body:    bodyJSON,
	}, nil
}
