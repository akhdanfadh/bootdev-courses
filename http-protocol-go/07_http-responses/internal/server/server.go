package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{
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
	headers := response.GetDefaultHeaders(0)
	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		log.Println("Error writing status line:", err)
		return
	}
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Println("Error writing headers:", err)
		return
	}
}
