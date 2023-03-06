// Package tcp provides a TCP transport
package tcp

import (
	"crypto/tls"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"kratos_t/transport_micro"
	"net"
	"time"
)

type tcpTransport struct {
	opts transport_micro.Options
}

type tcpTransportClient struct {
	dialOpts transport_micro.DialOptions
	conn     net.Conn
	timeout  time.Duration
}

type tcpTransportSocket struct {
	conn    net.Conn
	timeout time.Duration
}

type tcpTransportListener struct {
	listener net.Listener
	timeout  time.Duration
}

func (t *tcpTransportClient) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tcpTransportClient) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tcpTransportClient) Send(m *transport_micro.Message) error {
	return nil
}

func (t *tcpTransportClient) Recv(m *transport_micro.Message) error {
	return nil
}

func (t *tcpTransportClient) Close() error {
	return t.conn.Close()
}

func (t *tcpTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tcpTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tcpTransportSocket) Recv(m *transport_micro.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}

	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	by := make([]byte, 1024)
	n, err := t.conn.Read(by)
	if err != nil {
		return err
	}
	m.Body = by[:n]
	return nil
}

func (t *tcpTransportSocket) Send(m *transport_micro.Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	if _, err := t.conn.Write(m.Body); err != nil {
		return err
	}
	return nil
}

func (t *tcpTransportSocket) Close() error {
	return t.conn.Close()
}

func (t *tcpTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *tcpTransportListener) Close() error {
	return t.listener.Close()
}

func (t *tcpTransportListener) Accept(fn func(transport_micro.Socket)) error {
	var tempDelay time.Duration

	for {
		c, err := t.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 100 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Errorf("http: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		sock := &tcpTransportSocket{
			timeout: t.timeout,
			conn:    c,
		}

		go func() {
			// TODO: think of a better error response strategy
			defer func() {
				if r := recover(); r != nil {
					sock.Close()
				}
			}()

			fn(sock)
		}()
	}
}

func (t *tcpTransport) Dial(addr string, opts ...transport_micro.DialOption) (transport_micro.Client, error) {
	dopts := transport_micro.DialOptions{
		Timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(&dopts)
	}

	var conn net.Conn
	var err error

	// TODO: support dial option here rather than using internal config
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: dopts.Timeout}, "tcp", addr, config)
	} else {
		conn, err = net.DialTimeout("tcp", addr, dopts.Timeout)
	}

	if err != nil {
		return nil, err
	}

	return &tcpTransportClient{
		dialOpts: dopts,
		conn:     conn,
		timeout:  t.opts.Timeout,
	}, nil
}

func (t *tcpTransport) Listen(addr string, opts ...transport_micro.ListenOption) (transport_micro.Listener, error) {
	var options transport_micro.ListenOptions
	for _, o := range opts {
		o(&options)
	}

	var l net.Listener
	var err error

	l, err = net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	return &tcpTransportListener{
		timeout:  t.opts.Timeout,
		listener: l,
	}, nil
}

func (t *tcpTransport) Init(opts ...transport_micro.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *tcpTransport) Options() transport_micro.Options {
	return t.opts
}

func (t *tcpTransport) String() string {
	return "tcp"
}

func NewTransport(opts ...transport_micro.Option) transport_micro.Transport {
	var options transport_micro.Options
	for _, o := range opts {
		o(&options)
	}
	return &tcpTransport{opts: options}
}
