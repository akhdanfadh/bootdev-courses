package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Request represents an HTTP request, based on RFC 9112 Section 2.1.
type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}
type RequestLine struct {
	Method        string
	RequestTarget string
	HTTPVersion   string // initialism stylecheck HTTP inside Http
}

// RequestFromReader reads an HTTP request from the provided io.Reader.
func RequestFromReader(reader io.Reader) (*Request, error) {
	// Slurp the entire request into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Split the request into lines by CRLF, see RFC 9112 Section 2.1
	messageLines := strings.Split(string(data), "\r\n")
	// TODO: for now we only care about the first line, which should be the request line
	requestLine, err := parseRequestLine(messageLines[0])
	if err != nil {
		return nil, err
	}

	// Finally return the request
	return &Request{RequestLine: requestLine}, nil
}

// parseRequestLine parses the request line, based on RFC 9112 Section 3.
func parseRequestLine(line string) (RequestLine, error) {
	// Split on single space, expecting three parts: method, target, and version
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line: %s", line)
	}

	// Parse validity of each part
	method, err := parseMethod(parts[0])
	if err != nil {
		return RequestLine{}, err
	}
	target, err := parseRequestTarget(parts[1])
	if err != nil {
		return RequestLine{}, err
	}
	version, err := parseHTTPVersion(parts[2])
	if err != nil {
		return RequestLine{}, err
	}

	// Finally, return the parsed RequestLine
	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HTTPVersion:   version,
	}, nil
}

// parseMethod validates the HTTP method, based on RFC 9110 Section 3.1.
func parseMethod(s string) (string, error) {
	// TODO: only accpet any capitalized word
	if s == "" {
		return "", fmt.Errorf("invalid method: %s", s)
	}
	for _, r := range s {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
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
