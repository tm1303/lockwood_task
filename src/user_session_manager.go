package lockwood_task

import (
	"fmt"
	"lockwood_task/src/server"
	"sync"
)

type UserSessionManager struct {
	server server.PresenceServer
	mutex  sync.Mutex
	users  map[int]*UserSession
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

func (usm *UserSessionManager) UserConnects(request *server.ConnectionRequest, notifier server.UserNotifierChannel) {
	fmt.Printf("User Connecting: %v \n", request.UserId)
	userSession := NewUserSession(request.UserId, &request.Friends, notifier, usm)

	usm.mutex.Lock()
	defer usm.mutex.Unlock()

	usm.users[request.UserId] = userSession
	userSession.Notifier <- fmt.Sprintf("You have connected! (UserId: %v)", request.UserId)
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
