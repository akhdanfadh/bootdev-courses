package headers

import (
	"bytes"
	"fmt"
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
	key := strings.ToLower(string(newData[:colonIdx]))
	if key == "" {
		return 0, false, fmt.Errorf("invalid header key: %s", newData)
	}
	if len(key) != len(strings.TrimRightFunc(key, unicode.IsSpace)) {
		return 0, false, fmt.Errorf("invalid header key with whitespace: %s", newData)
	}
	key = strings.TrimSpace(key) // trim leading whitespace
	// value could be empty
	value := string(bytes.TrimSpace(newData[colonIdx+1:]))

	h[key] = value // dereference pointer to access map
	return CLRFIdx + 2, false, nil
}
