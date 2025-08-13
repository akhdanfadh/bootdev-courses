package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/headers"
	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/request"
	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/response"
)

const (
	chunkSize   = 1024 // max bytes per chunk
	httpbinBase = "https://httpbin.org/"
)

func httpbinHandler(w *response.Writer, req *request.Request) {
	// Parse original request target
	origTarget := req.RequestLine.RequestTarget
	if !strings.HasPrefix(origTarget, "/httpbin/") {
		httpbinHandlerError(w, response.StatusBadRequest, "Invalid request target")
		return
	}

	// Redirect to the actual URL and get the response
	url := httpbinBase + strings.TrimPrefix(origTarget, "/httpbin/")
	resp, err := http.Get(url)
	if err != nil {
		httpbinHandlerError(w, response.StatusInternalServerError, fmt.Sprintf("Failed to fetch %s: %v", url, err))
		return
	}
	defer resp.Body.Close()

	// Write response to user from upstream response
	w.WriteStatusLine(response.StatusOK)
	h := headers.NewHeaders()
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(h)

	buffer := make([]byte, chunkSize)
	fullBody := make([]byte, 0)
	for {
		n, err := resp.Body.Read(buffer)
		if n == 0 && err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading response body: %v", err)
			break
		}
		if _, err := w.WriteChunkedBody(buffer[:n]); err != nil {
			log.Printf("Error writing chunked body: %v", err)
			break
		}
		fullBody = append(fullBody, buffer[:n]...)
	}
	if _, err := w.WriteChunkedBodyDone(); err != nil {
		log.Printf("Error writing chunked body: %v", err)
	}

	// Validate trailer
	t := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	t.Set("X-Content-SHA256", sha256)
	t.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
	w.WriteTrailer(t)
}

func httpbinHandlerError(w *response.Writer, statusCode response.StatusCode, message string) {
	w.WriteStatusLine(statusCode)
	body := []byte(message)
	w.WriteHeaders(response.GetDefaultHeaders(len(body)))
	w.WriteBody(body)
}
