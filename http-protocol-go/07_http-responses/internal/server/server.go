package server

import (
	"fmt"
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

type Handler func(w *response.Writer, req *request.Request)

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
	w := response.NewWriter(conn)

	// Parse the request from connection
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := fmt.Appendf(nil, "Error parsing request: %v", err)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req) // handle if no error
}
