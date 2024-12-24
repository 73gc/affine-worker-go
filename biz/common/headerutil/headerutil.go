package headerutil

import (
	"regexp"

	"github.com/cloudwego/hertz/pkg/protocol"
)

func CloneHeaders(rawHeader *protocol.RequestHeader) (map[string]string, error) {
	headerFilters := []*regexp.Regexp{
		regexp.MustCompile(`(?i)^Sec-`),
		regexp.MustCompile(`(?i)^Accept`),
		regexp.MustCompile(`(?i)^User-Agent$`),
	}
	headers := make(map[string]string)
	f := func(key, value []byte) {
		headers[string(key)] = string(value)
	}
	rawHeader.VisitAll(f)
	resHeaders := make(map[string]string)
	for key, value := range headers {
		for _, filter := range headerFilters {
			if filter.MatchString(key) {
				resHeaders[key] = value
			}
		}
	}
	return resHeaders, nil
}
