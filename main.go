package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos_t/server"
	"kratos_t/transport_micro"
	"kratos_t/transport_micro/udp"
)

func main() {
	err := kratos.New(
		kratos.Server(
			server.NewServer(
				server.WithTransport(udp.NewTransport()),
				server.WithRootHandler(Echo),
			),
		),
	).Run()
	if err != nil {
		panic(err)
	}
}

func Echo(sock transport_micro.Socket, msg []byte) {
	err := sock.Send(&transport_micro.Message{Body: msg})
	if err != nil {
		log.Error(err)
	}
}
