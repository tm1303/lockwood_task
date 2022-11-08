package server

import (
	"fmt"
	"net"
)

type PresenceServer interface {
	Listen(handler func(*ConnectionRequest))
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

func (s *UdpServer) Listen(handler func(*ConnectionRequest)) {

	buffer := make([]byte, 1024)
	for {
		n, _, _ := s.connection.ReadFromUDP(buffer)
		bits := buffer[0:n]

		conRequest := NewConnectionRequest(bits)
		fmt.Printf("conRequest, UserId: %v \n", conRequest.UserId)
	}
}
