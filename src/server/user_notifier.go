package server

import (
	"net"
)

type UserNotifierChannel chan string

func NewUdpNotifier(server PresenceServer, addr *net.UDPAddr) UserNotifierChannel {

	messageChan := make(chan string)
	go func() {
		for {
			message := <-messageChan
			server.WriteUdp(message, addr)
		}
	}()
	return messageChan
}

func NewTcpNotifier(server PresenceServer, vars any) UserNotifierChannel {
	panic("Not yet implemented.")
	// but when it is we should have a nice generic way to notify users connected to either flavour of server
}
