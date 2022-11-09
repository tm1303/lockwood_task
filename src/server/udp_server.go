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
