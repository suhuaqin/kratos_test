package server

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"io"
	"kratos_t/transport_micro"
)

var _ transport.Server = (*Server)(nil)

type Server struct {
	rootHandler RootHandler
	tr          transport_micro.Transport
}

func (s *Server) init() error {
	return nil
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		rootHandler: nil,
		tr:          nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.init()
	return s
}

func (s *Server) Start(ctx context.Context) error {
	ln, err := s.tr.Listen("127.0.0.1:6666")
	if err != nil {
		return err
	}

	for {
		err := ln.Accept(s.ServerConn)
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *Server) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ServerConn(socket transport_micro.Socket) {
	defer socket.Close()
	msg := &transport_micro.Message{}
	for {
		err := socket.Recv(msg)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Error(err)
			}
			return
		}

		s.rootHandler(socket, msg.Body)
	}
}
