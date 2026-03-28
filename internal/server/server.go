package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net"
)

type Server struct {
	handler response.Handler
	closed  bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	// out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	// conn.Write(out)
	// conn.Close()
	defer conn.Close()

	netConn := conn.(net.Conn)
	s.handle(netConn)
}

func runServer(s *Server, listner net.Listener) {
	for {
		conn, err := listner.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port uint16, handler response.Handler) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{handler: handler, closed: false}
	go runServer(server, listner)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}

func (s *Server) listen() {

}

func (s *Server) handle(conn net.Conn) {
	responseWriter := response.NewWriter(conn)

	headers := response.GetDefaultHeaders(0)
	// IMPORTANT: the whole parsing of the request lies here
	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(headers)
		return
	}

	s.handler(responseWriter, r)
}
