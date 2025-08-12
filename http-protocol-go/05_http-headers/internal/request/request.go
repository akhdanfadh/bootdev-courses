package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/headers"
)

// Request represents an HTTP request, based on RFC 9112 Section 2.1.
type (
	Request struct {
		state       parseState
		RequestLine RequestLine
		Headers     headers.Headers
		Body        []byte
	}
	RequestLine struct {
		Method        string
		RequestTarget string
		HTTPVersion   string
	}
	parseState int
)

const bufferSize = 8
const (
	isDone parseState = iota
	isRequestLine
	isHeaders
	isBody
)

// RequestFromReader reads an HTTP request from the provided io.Reader.
func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize) // buffer to read data into
	readToIndex := 0                   // keep track how much data we've read
	request := &Request{
		state:   isRequestLine,        // initialize the request parse state
		Headers: headers.NewHeaders(), // initialize headers
	}
	for request.state != isDone {
		// Grow the buffer if full
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, cap(buffer)*2) // double the capacity
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		// Read into buffer starting at readToIndex
		bytesRead, err := reader.Read(buffer[readToIndex:])
		if bytesRead == 0 && err != nil {
			if err == io.EOF {
				request.state = isDone
				break
			}
			return nil, err
		}
		readToIndex += bytesRead

		// Parse data we've read so far
		bytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Remove data that has been parsed to keep buffer small
		if bytesParsed > 0 {
			newBuffer := make([]byte, len(buffer)-bytesParsed, cap(buffer))
			copy(newBuffer, buffer[bytesParsed:readToIndex])
			buffer = newBuffer
			readToIndex -= bytesParsed // not forgetting this
		}
	}

	// Finally return the request
	return request, nil
}

// parse processes the data read so far and based on the current state
func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case isDone:
		return 0, fmt.Errorf("request is already done")
	case isRequestLine:
		requestLine, bytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			return 0, nil
		} // not enough data to parse request line yet
		r.RequestLine = *requestLine
		r.state = isHeaders // move to the next state
		return bytesParsed, nil
	case isHeaders:
		bytesParsed, doneParsing, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			return 0, nil
		} // not enough data to parse headers yet
		if doneParsing {
			r.state = isDone
		} // TODO: for now we only parse request line and headers
		return bytesParsed, nil
	case isBody:
		return 0, fmt.Errorf("parsing body is not implemented yet")
	default:
		return 0, fmt.Errorf("unknown parse state: %d", r.state)
	}
}

// parseRequestLine parses the request line, based on RFC 9112 Section 3.
func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// Only parse if there is CRLF in the data
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil // no CRLF found, nothing to parse, need more data
	}

	// Split on single space, expecting three parts: method, target, and version
	line := string(data[:idx])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("invalid request line: %s", line)
	}

	// Parse validity of each part
	method, err := parseMethod(parts[0])
	if err != nil {
		return nil, 0, err
	}
	target, err := parseRequestTarget(parts[1])
	if err != nil {
		return nil, 0, err
	}
	version, err := parseHTTPVersion(parts[2])
	if err != nil {
		return nil, 0, err
	}

	// Finally, return the parsed RequestLine
	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HTTPVersion:   version,
	}, idx + 2, nil // 2 accounting CRLF bytes
}

// parseMethod validates the HTTP method, based on RFC 9110 Section 3.1.
func parseMethod(s string) (string, error) {
	// TODO: only accept any capitalized word
	if s == "" {
		return "", fmt.Errorf("invalid method: %s", s)
	}
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsLetter(r) {
			return "", fmt.Errorf("invalid method: %s", s)
		}
	}
	return s, nil
}

// parseRequestTarget validates the request target, based on RFC 9112 Section 3.2.
func parseRequestTarget(s string) (string, error) {
	// TODO: only accept non-empty string
	if s == "" {
		return "", fmt.Errorf("invalid request target: %s", s)
	}
	return s, nil
}

// parseHTTPVersion validates the HTTP version, based on RFC 9112 Section 2.3.
func parseHTTPVersion(s string) (string, error) {
	// TODO: only accept HTTP/1.1
	if s != "HTTP/1.1" {
		return "", fmt.Errorf("invalid HTTP version: %s", s)
	}
	return "1.1", nil
}
