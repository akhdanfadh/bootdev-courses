package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	// Only parse if there is CRLF in the data
	CLRFIdx := bytes.Index(data, []byte("\r\n"))
	if CLRFIdx == -1 {
		return 0, false, nil // no CRLF found, nothing to parse, need more data
	}
	if CLRFIdx == 0 {
		return 2, true, nil // end of headers, CRLF found at the start
	}

	// Parse header, based on RFC 9112 Section 5
	newData := bytes.TrimSpace(data[:CLRFIdx])
	colonIdx := bytes.Index(newData, []byte(":"))
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("invalid header format: %s", newData)
	}
	// key is case-insensitive and no whitespace allowed before colon
	key, err := parseHeaderKey(newData[:colonIdx])
	if err != nil {
		return 0, false, err
	}
	// value could be empty
	value := string(bytes.TrimSpace(newData[colonIdx+1:]))

	h.Set(key, value)
	return CLRFIdx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	if val, ok := h[key]; ok { // if key exists
		h[key] = val + ", " + value
	} else {
		h[key] = value
	}
}

var headerKeySymbols = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func parseHeaderKey(data []byte) (string, error) {
	key := string(data)
	if key == "" {
		return "", fmt.Errorf("invalid empty header key")
	}
	if key != strings.TrimRightFunc(key, unicode.IsSpace) {
		return "", fmt.Errorf("invalid whitespace in header key: '%s'", data)
	}
	key = strings.ToLower(strings.TrimSpace(key)) // trim leading whitespace and case-insensitive
	for _, r := range key {
		if (r < 'a' || r > 'z') &&
			(r < '0' || r > '9') &&
			!slices.Contains(headerKeySymbols, byte(r)) {
			return "", fmt.Errorf("invalid characters in header key: '%s'", data)
		}
	}
	return key, nil
}
