package network

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Exec execute HTTP requests according to following parameters.
func Exec(method, url string, body []byte, headers map[string]string, cookie *http.Cookie) ([]*http.Cookie, []byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Print("[ERROR] Request error", err)
		return nil, nil, err
	}

	for key, header := range headers {
		req.Header.Add(key, header)
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Print("[ERROR] Request error", err)
		return nil, nil, err
	}
	defer res.Body.Close()

	bodyJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("[ERROR] Request error", err)
		fmt.Println(res.StatusCode, string(bodyJSON))
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, nil, errors.New("bad request")
	}

	return res.Cookies(), bodyJSON, nil
}
