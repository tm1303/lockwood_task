package lockwood_task

import (
	"fmt"
	"lockwood_task/src/server"
	"net"
)

type ConnectionManager struct {
	users  map[UserId]*UserSession
	server server.PresenceServer
}

func NewConnectionManager(server server.PresenceServer) *ConnectionManager {
	return &ConnectionManager{
		users:  make(map[UserId]*UserSession, 0),
		server: server,
	}
}

func (cm *ConnectionManager) Start() {
	cm.server.Listen(cm.UserConnects)
}

func (cm *ConnectionManager) UserConnects(request *server.ConnectionRequest, addr *net.UDPAddr) {
	fmt.Printf("User Connecting: %v \n", request.UserId)

	

	cm.server.Write(fmt.Sprintf("You have connected! (UserId: %v)", request.UserId), addr)
}
