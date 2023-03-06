package udp

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"kratos_t/transport_micro"
	"net"
)

func (u *udpListener) Addr() string {
	return u.l.LocalAddr().String()
}

func (u *udpListener) Close() error {
	return u.l.Close()
}

func (u *udpListener) Accept(fn func(transport_micro.Socket)) error {
	by := make([]byte, 1024*1024)
	go func() {
		for {
			select {
			case msg := <-u.send:
				_, err := u.l.WriteToUDP(msg.body, msg.addr)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()

	for {
		n, cAddr, err := u.l.ReadFromUDP(by)
		if err != nil {
			return err
		}

		sock, isNew := u.addSession(cAddr)
		if isNew {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						sock.Close()
					}
				}()

				fn(sock)
			}()
		}
		// TODO: 复用
		data := make([]byte, n)
		copy(data, by)
		fmt.Println(string(data))
		sock.rcv <- &transport_micro.Message{
			Body: data,
		}
	}
}

func (u *udpListener) addSession(addr *net.UDPAddr) (*udpSocket, bool) {
	sockID := sockID(addr.String())
	u.session.RLock()
	sock := u.session.session[sockID]
	u.session.RUnlock()
	if sock != nil {
		return sock, false
	}
	u.session.Lock()
	defer u.session.Unlock()
	u.session.session[sockID] = &udpSocket{
		rcv:    make(chan *transport_micro.Message),
		send:   u.send,
		listen: u.l,
		addr:   addr,
	}
	return u.session.session[sockID], true
}

func (u *udpListener) delSession() {

}
