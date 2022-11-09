package lockwood_task

import (
	"fmt"
	"lockwood_task/src/server"
	"net"
	"sync"
)

type UserSessionManager struct {
	mutex  *sync.Mutex
	users  map[int]*UserSession
	server server.PresenceServer
}

func NewUserSessionManager(server server.PresenceServer) *UserSessionManager {
	return &UserSessionManager{
		users:  make(map[int]*UserSession, 0),
		server: server,
	}
}

func (usm *UserSessionManager) Start() {
	usm.server.Listen(usm.UserConnects)
}

func (usm *UserSessionManager) UserConnects(request *server.ConnectionRequest, addr *net.UDPAddr) {
	fmt.Printf("User Connecting: %v \n", request.UserId)
	usm.server.Write(fmt.Sprintf("You have connected! (UserId: %v)", request.UserId), addr)

	session := NewUserSession(request.UserId, &request.Friends, &Connection{Udp: addr}, usm)

	usm.mutex.Lock()
	defer usm.mutex.Unlock()

	usm.users[request.UserId] = session
}

func (usm *UserSessionManager) GetConnectedUser(userId int) (userSession *UserSession, found bool) {

	usm.mutex.Lock()
	defer usm.mutex.Unlock()

	if user, found := usm.users[userId]; found {
		return user, true
	} else {
		return OfflineUser, false
	}
}
