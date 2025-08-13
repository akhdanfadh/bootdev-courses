package main

import (
	"fmt"
	"os"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/headers"
	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/request"
	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/response"
)

const videoPath = "assets/vim.mp4"

func videoHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget != "/video" {
		httpbinHandlerError(w, response.StatusBadRequest, "Invalid request target, expected /video")
		return
	}

	data, err := os.ReadFile(videoPath)
	if err != nil {
		httpbinHandlerError(w, response.StatusInternalServerError, "Could not read video file")
		return
	}

	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	h.Set("Connection", "close")
	h.Set("Content-Type", "video/mp4")
	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(h)
	w.WriteBody(data)
}
