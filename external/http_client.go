package external

import "net/http"

var client *http.Client

func getHttpClient() *http.Client {
	if client == nil {
		client = &http.Client{}
	}
	return client
}
