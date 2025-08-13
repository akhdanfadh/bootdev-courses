package response

import (
	"fmt"
	"io"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/headers"
)

type Writer struct {
	writer io.Writer
	state  writerState
}
type writerState int

const (
	isStatusLine writerState = iota
	isHeaders
	isBody
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  isStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != isStatusLine {
		return fmt.Errorf("cannot write status line in state %v", w.state)
	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText[statusCode])
	_, err := w.writer.Write([]byte(statusLine))
	if err == nil {
		w.state = isHeaders
	}
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != isHeaders {
		return fmt.Errorf("cannot write headers in state %v", w.state)
	}
	for key, value := range headers {
		if _, err := fmt.Fprintf(w.writer, "%s: %s\r\n", key, value); err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n")) // end of headers
	if err == nil {
		w.state = isBody
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != isBody {
		return 0, fmt.Errorf("cannot write body in state %v", w.state)
	}
	return w.writer.Write(p)
}
