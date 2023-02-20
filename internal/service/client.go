package service

import (
	"os"
	"time"

	"github.com/imroc/req/v3"
)

const (
	DumpRequest = "DUMP_REQUEST"
)

var reqClient *req.Client

// NewClient returns a req client. Singleton pattern is used to create the client
// only once. Default headers are added for Accept and Accept-Language schemas.
// To use this client, call the R() method to create a unique request.
func NewClient() *req.Client {
	if reqClient == nil {
		client := req.C(). // Use C() to create a client and set with chainable client settings.
					SetUserAgent("etna-notification-bot").
					SetTimeout(15 * time.Second).
					EnableDebugLog().
					SetCookieJar(nil). // Disable cookie storage between requests
					SetCommonHeaders(map[string]string{
				"Accept":          "application/json",
				"Accept-Language": "fr-FR",
			})
		reqClient = client

		if os.Getenv(DumpRequest) == "true" {
			client.DevMode()
		}
	}

	return reqClient
}
