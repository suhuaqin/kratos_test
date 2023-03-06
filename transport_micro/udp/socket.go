package udp

import (
	"errors"
	"kratos_t/transport_micro"
	"net"
)

func (u *udpSocket) Local() string {
	return u.listen.LocalAddr().String()
}

func (u *udpSocket) Remote() string {
	return u.addr.String()
}

func (u *udpSocket) Recv(m *transport_micro.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}
	*m = *<-u.rcv
	// TODO: EOF
	return nil
}

func (u *udpSocket) Send(m *transport_micro.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}
	u.send <- &udpMessage{
		addr: u.addr,
		body: m.Body,
	}
	return nil
}

func (u *udpSocket) Close() error {
	return nil
}

type udpMessage struct {
	addr *net.UDPAddr
	body []byte
}
