package udp

import (
	"kratos_t/transport_micro"
	"net"
	"sync"
)

type udpTransport struct {
	opts transport_micro.Options
}

type udpListener struct {
	opts    transport_micro.ListenOptions
	session *session
	send    chan *udpMessage
	l       *net.UDPConn
}

type utpClient struct {
	dialOpts transport_micro.DialOptions
	conn     *net.UDPConn
}

type udpSocket struct {
	rcv    chan *transport_micro.Message
	send   chan *udpMessage
	listen *net.UDPConn
	addr   *net.UDPAddr
}

func NewTransport(opts ...transport_micro.Option) transport_micro.Transport {
	var options transport_micro.Options
	for _, o := range opts {
		o(&options)
	}
	return &udpTransport{opts: options}
}

type sockID string

type session struct {
	session map[sockID]*udpSocket
	sync.RWMutex
}

func newSession() *session {
	return &session{
		session: make(map[sockID]*udpSocket),
	}
}
