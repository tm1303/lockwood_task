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

func (usm *UserSessionManager) UserConnects(request *server.LogOnRequest, notifier server.UserNotifierChannel) {

	userSession, found := usm.GetConnectedUser(request.UserId) // todo: could poss check udp address to ensure is same user client??
	if !found {
		usm.AddUser(request, userSession, notifier)
	} else {
		userSession.ResetTimeout()
		userSession.Notifier = notifier //because their udp address changed
	}
}

func (usm *UserSessionManager) AddUser(request *server.LogOnRequest, userSession *UserSession, notifier server.UserNotifierChannel) {
	fmt.Printf("User Connecting: %v \n", request.UserId)
	userSession = NewUserSession(request.UserId, &request.Friends, notifier, usm)

	usm.mutex.Lock()
	defer usm.mutex.Unlock()

	usm.users[request.UserId] = userSession
	userSession.Notifier <- fmt.Sprintf("You have connected! (UserId: %v)", request.UserId)

	userSession.ResetTimeout()
	go func() {
		for {
			user, found := usm.GetConnectedUser(request.UserId)
			if !found{
				break
			}
			// if the users SessionTimeout expires before it is reset we must kill the user!
			<-user.SessionTimeout.C
			usm.RemoveUser(userSession)
		}
	}()
}

func (usm *UserSessionManager) RemoveUser(userSession *UserSession) {
	fmt.Printf("User removed: %v \n", userSession.userId)
	usm.mutex.Lock()
	defer usm.mutex.Unlock()
	delete(usm.users, userSession.userId)
}

func (usm *UserSessionManager) GetConnectedUser(userId int) (userSession *UserSession, found bool) {

	usm.mutex.Lock()
	defer usm.mutex.Unlock()

	if user, found := usm.users[userId]; found {
		return user, user != nil
	} else {
		return nil, false
	}
}

func (usm *UserSessionManager) VerifyConnectedUser(user *UserSession) (found bool) {

	usm.mutex.Lock()
	defer usm.mutex.Unlock()

	if foundUser, found := usm.users[user.userId]; !found {
		return false
	} else {
		return foundUser == user
	}
}
