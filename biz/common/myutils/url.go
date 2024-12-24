package myutils

import (
	"affine-worker-go/biz/common/tld"
	"errors"
	"net/url"
	"strings"
)

func FixURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, nil
	}
	fullURL := rawURL
	if !strings.HasPrefix(fullURL, "http:") && !strings.HasPrefix(fullURL, "https:") {
		fullURL = "http:" + fullURL
	}
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	parsedRawURL, err := tld.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	var fullDomain string
	if parsedRawURL.Subdomain != "" {
		fullDomain = parsedRawURL.Subdomain + "." + parsedRawURL.Domain
	} else {
		fullDomain = parsedRawURL.Domain
	}
	if fullDomain == parsedURL.Hostname() {
		return parsedURL, nil
	}
	return nil, errors.New("invalid url")
}
