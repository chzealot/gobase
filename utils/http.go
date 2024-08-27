package utils

import (
	"fmt"
	"github.com/chzealot/gobase/constants"
	"github.com/chzealot/gobase/logger"
	"io"
	"net/http"
	"strings"
)

type getHeaderInterface interface {
	GetHeader(key string) string
}

func DumpHttpRequest(r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		body = make([]byte, 0)
		return
	}
	headers := []string{}
	for name, values := range r.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", name, value))
		}
	}
	logger.Infof("%s %s\nHost: %s\n%s\n\n%s",
		r.Method, r.URL.String(),
		r.Host,
		strings.Join(headers, "\n"),
		string(body))
}

func GetHttpProto(c getHeaderInterface) string {
	scheme := "http"
	if c.GetHeader(constants.HeaderForwardedProto) == "https" {
		scheme = "https"
	}
	return scheme
}
