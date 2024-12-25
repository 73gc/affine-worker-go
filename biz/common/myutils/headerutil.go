package myutils

import (
	"regexp"

	"github.com/cloudwego/hertz/pkg/protocol"
)

func CloneHeaders(rawHeader *protocol.RequestHeader) map[string]string {
	headerFilters := []*regexp.Regexp{
		regexp.MustCompile(`(?i)^Sec-`),
		regexp.MustCompile(`(?i)^Accept`),
		regexp.MustCompile(`(?i)^User-Agent$`),
	}
	resHeaders := make(map[string]string)
	f := func(key, value []byte) {
		for _, filter := range headerFilters {
			if filter.MatchString(string(key)) {
				resHeaders[string(key)] = string(value)
			}
		}
	}
	rawHeader.VisitAll(f)
	return resHeaders
}
