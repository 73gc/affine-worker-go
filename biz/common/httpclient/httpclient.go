package httpclient

import (
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: time.Minute,
}

func Do(req *http.Request) (*http.Response, error) {
	return Client.Do(req)
}
