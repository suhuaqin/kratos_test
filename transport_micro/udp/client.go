package udp

import "kratos_t/transport_micro"

func (u *utpClient) Local() string {
	return u.conn.LocalAddr().String()
}

func (u *utpClient) Remote() string {
	return u.conn.RemoteAddr().String()
}

func (u *utpClient) Send(m *transport_micro.Message) error {
	return nil
}

func (u *utpClient) Recv(m *transport_micro.Message) error {
	return nil
}

func (u *utpClient) Close() error {
	return u.conn.Close()
}
