package tld

import (
	"net"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type URL struct {
	*url.URL
	Domain              string // TLD + SLD
	Subdomain           string
	PublicSuffix        string
	DomainWithoutSuffix string
	Iccan               bool // Does TLD come from ICANN part of the list
	IsIp                bool
	IsPrivate           bool // Does TLD come from Private part of the list
}

func Parse(rawURL string) (*URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	hostname := strings.ToLower(parsedURL.Hostname())
	parsedIP := net.ParseIP(hostname)
	if parsedIP != nil {
		return &URL{
			IsIp: true,
			URL:  parsedURL,
		}, nil
	}
	suffix, iccan := publicsuffix.PublicSuffix(hostname)
	if !iccan {
		return &URL{
			PublicSuffix: suffix,
			Iccan:        iccan,
			IsIp:         false,
			IsPrivate:    !iccan,
			URL:          parsedURL,
		}, nil
	}
	etld1, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return nil, err
	}
	subdomain := strings.TrimSuffix(hostname, etld1)
	subdomain = strings.TrimSuffix(subdomain, ".")
	domainWithoutSuffix := strings.TrimSuffix(etld1, suffix)
	domainWithoutSuffix = strings.TrimSuffix(domainWithoutSuffix, ".")
	return &URL{
		Domain:              etld1,
		Subdomain:           subdomain,
		PublicSuffix:        suffix,
		DomainWithoutSuffix: domainWithoutSuffix,
		Iccan:               iccan,
		IsIp:                false,
		IsPrivate:           !iccan,
		URL:                 parsedURL,
	}, nil
}
