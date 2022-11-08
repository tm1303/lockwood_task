package lockwood_task

import "lockwood_task/src/server"

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

func (cm *ConnectionManager) UserConnects(request *server.ConnectionRequest) {
	//to do: handle a user connecting
}
