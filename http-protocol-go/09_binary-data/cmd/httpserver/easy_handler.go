package main

import (
	"fmt"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/headers"
	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/request"
	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/response"
)

func easyHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
		writeDefaultEasyHandler(w, body)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		w.WriteStatusLine(response.StatusInternalServerError)
		body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
		writeDefaultEasyHandler(w, body)
		return
	}

	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	writeDefaultEasyHandler(w, body)
}

func writeDefaultEasyHandler(w *response.Writer, body []byte) {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
