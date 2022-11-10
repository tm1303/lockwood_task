package server

import (
	"fmt"
	"net"
)

type PresenceServer interface {
	Listen(handler func(*LogOnRequest, UserNotifierChannel))
	WriteUdp(message string, addr *net.UDPAddr)
	WriteTcp(message string, vars any) // todo: support tcp
}

type UdpServer struct {
	connection *net.UDPConn
}

func NewUdpServer(port string) *UdpServer {
	s, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return nil
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &UdpServer{connection: connection}
}

func (s *UdpServer) Listen(handler func(*LogOnRequest, UserNotifierChannel)) {
	fmt.Println("Presence Server listening for user logins")
	buffer := make([]byte, 1024)
	for {
		n, addr, _ := s.connection.ReadFromUDP(buffer)

		conRequest := ParseLogOnRequest(buffer[0:n])
		// oh no, I've made a mistake with no time to fix it
		// I'm creating a new Notifier every time a client pings us
		// I should only do this when the addr for a userId changes
		// TODO: fix it, but not today.
		notifier := NewUdpNotifier(s, addr)
		handler(conRequest, notifier)
	}
}

func (s *UdpServer) WriteUdp(message string, addr *net.UDPAddr) {

	data := []byte(message)
	_, err := s.connection.WriteToUDP(data, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *UdpServer) WriteTcp(message string, vars any) {
	panic("UDP Server doesn't support TCP")
}
