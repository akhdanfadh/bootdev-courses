package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/request"
	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/response"
)

type Server struct {
	handler  Handler
	listener net.Listener
	isClosed atomic.Bool
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

// Write writes the handler error response to the provided writer.
func (he *HandlerError) Write(w io.Writer) error {
	if err := response.WriteStatusLine(w, he.StatusCode); err != nil {
		return err
	}
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	if err := response.WriteHeaders(w, headers); err != nil {
		return err
	}
	if _, err := w.Write(messageBytes); err != nil {
		return err
	}
	return nil
}

func Serve(port int, handler Handler) (*Server, error) {
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		handler:  handler,
		listener: l,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true) // mark server as closed
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return // exit silently if server is closed
			}
			log.Println("Error accepting connection:", err)
			continue // continue accpting new connections even if one fails
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close() // ensure connection closed after handling

	// Parse the request from connection
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}

	// Handle the request using the provided handler
	buf := bytes.NewBuffer(nil) // to hold the response body
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}

	// If all good, write the response
	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		log.Println("Error writing status line:", err)
		return
	}
	data := buf.Bytes()
	headers := response.GetDefaultHeaders(len(data))
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Println("Error writing headers:", err)
		return
	}
	if _, err := conn.Write(data); err != nil {
		log.Println("Error writing response body:", err)
		return
	}
}
