package lockwood_task

import "net"

// I wonder if we can make this nice and generic, swappable TCP/UDP
type Connection struct {
	Udp *net.UDPAddr
	Tcp *net.Conn
}
