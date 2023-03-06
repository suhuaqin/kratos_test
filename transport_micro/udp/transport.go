package udp

import (
	"kratos_t/transport_micro"
	"net"
)

func (u *udpTransport) Dial(addr string, opts ...transport_micro.DialOption) (transport_micro.Client, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	dopts := transport_micro.DialOptions{}

	for _, opt := range opts {
		opt(&dopts)
	}

	c, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &utpClient{
		dialOpts: dopts,
		conn:     c,
	}, nil
}

func (u *udpTransport) Listen(addr string, opts ...transport_micro.ListenOption) (transport_micro.Listener, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	var options transport_micro.ListenOptions
	for _, o := range opts {
		o(&options)
	}

	l, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &udpListener{
		opts:    options,
		session: newSession(),
		send:    make(chan *udpMessage, 256),
		l:       l,
	}, nil
}

func (u *udpTransport) Init(opts ...transport_micro.Option) error {
	for _, o := range opts {
		o(&u.opts)
	}
	return nil
}

func (u *udpTransport) Options() transport_micro.Options {
	return u.opts
}

func (u *udpTransport) String() string {
	return "udp"
}
