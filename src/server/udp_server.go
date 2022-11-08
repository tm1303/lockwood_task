package server

import (
	"fmt"
	"net"
)

type PresenceServer interface {
	Listen(handler func(*ConnectionRequest, *net.UDPAddr))
	Write(message string, addr *net.UDPAddr) 
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

func (s *UdpServer) Listen(handler func(*ConnectionRequest, *net.UDPAddr)) {

	buffer := make([]byte, 1024)
	for {
		n, addr, _ := s.connection.ReadFromUDP(buffer)

		conRequest := NewConnectionRequest(buffer[0:n])
		handler(conRequest, addr)
	}
}

func (s *UdpServer) Write(message string, addr *net.UDPAddr) {

	data := []byte(message)
	fmt.Printf("data: %s\n", string(data))
	_, err := s.connection.WriteToUDP(data, addr)
	if err != nil {
			fmt.Println(err)
			return
	}
}
