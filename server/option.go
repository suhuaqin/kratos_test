package server

import (
	"kratos_t/transport_micro"
)

type ServerOption func(*Server)

type RootHandler func(sock transport_micro.Socket, msg []byte)

func WithTransport(transport transport_micro.Transport) ServerOption {
	return func(s *Server) {
		s.tr = transport
	}
}

func WithRootHandler(rootHandler RootHandler) ServerOption {
	return func(s *Server) {
		s.rootHandler = rootHandler
	}
}
