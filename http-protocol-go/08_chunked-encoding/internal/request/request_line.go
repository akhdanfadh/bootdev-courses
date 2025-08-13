package request

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HTTPVersion   string
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
