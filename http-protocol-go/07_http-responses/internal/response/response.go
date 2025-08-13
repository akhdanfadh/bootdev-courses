package response

import (
	"fmt"
	"io"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

var statusText = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText[statusCode])
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", key, value); err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n")) // end of headers
	return err
}
