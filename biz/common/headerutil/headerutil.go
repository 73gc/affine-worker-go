package headerutil

import (
	"regexp"

	"github.com/cloudwego/hertz/pkg/protocol"
)

func CloneHeaders(rawHeader *protocol.RequestHeader) (map[string]interface{}, error) {
	headerFilters := []*regexp.Regexp{
		regexp.MustCompile(`(?i)^Sec-`),
		regexp.MustCompile(`(?i)^Accept`),
		regexp.MustCompile(`(?i)^User-Agent$`),
	}
	headers := make(map[string]interface{})
	f := func(key, value []byte) {
		headers[string(key)] = string(value)
	}
	rawHeader.VisitAll(f)
	resHeaders := make(map[string]interface{})
	for key, value := range headers {
		for _, filter := range headerFilters {
			if filter.MatchString(key) {
				resHeaders[key] = value
			}
		}
	}
	return resHeaders, nil
}
