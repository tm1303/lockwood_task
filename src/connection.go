package lockwood_task

import "net"

type Connection struct {
	Udp *net.UDPAddr
	Tcp *net.Conn
}
