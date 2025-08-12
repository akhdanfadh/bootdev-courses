package request

import (
	"fmt"
	"io"
	"strconv"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/headers"
)

// Request represents an HTTP request, based on RFC 9112 Section 2.1.
type Request struct {
	state       parseState
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
}
type parseState int

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
		Body:    make([]byte, 0),      // initialize body
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
				if request.state != isDone {
					return nil, fmt.Errorf("incomplete request, currently in state: %v", request.state)
				}
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
			copy(buffer, buffer[bytesParsed:])
			readToIndex -= bytesParsed // not forgetting this
		}
	}

	// Finally return the request
	return request, nil
}

// parse process the data we have so far (previous data + a read from io.Reader).
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != isDone {
		// - Our parsing function is designed to handle one 'line' at at time,
		//   by checking if there is a CRLF sequence in the data.
		// - However, there could be multiple valid parseable parts in given data,
		//   e.g., `\r\n\r\n` just before request body.
		// - This parts could lost if we don't parse it again, e.g., if ONLY one
		//   parse call is done every read AND next read is EOF.
		// - Thus we need to handle multiple parse calls gracefully here.
		bytesParsed, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			break // need more data to read
		}
		totalBytesParsed += bytesParsed
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
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
		if doneParsing {
			r.state = isBody // move to the next state
		}
		return bytesParsed, nil
	case isBody:
		// Validate content-length header
		val, found := r.Headers.Get("content-length")
		if !found {
			r.state = isDone      // TODO: final state, we simply assume no body
			return len(data), nil // report data length even since we won't parsed it
		}
		num, err := strconv.Atoi(val)
		if err != nil || num < 0 {
			return 0, fmt.Errorf("invalid content-length header: %s", val)
		}
		// Append data to the body
		r.Body = append(r.Body, data...)
		if len(r.Body) > num {
			return 0, fmt.Errorf("body exceeds content-length: %d > %d", len(r.Body), num)
		}
		if len(r.Body) == num {
			r.state = isDone // move to the final state
		}
		return len(data), nil
	default:
		return 0, fmt.Errorf("unknown parse state: %d", r.state)
	}
}
