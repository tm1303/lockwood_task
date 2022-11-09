package lockwood_task

import (
	"fmt"
	"lockwood_task/src/server"
	"net"
	"sync"
)

type ConnectionManager struct {
	mutex  *sync.Mutex
	users  map[int]*UserSession
	server server.PresenceServer
}

func NewConnectionManager(server server.PresenceServer) *ConnectionManager {
	return &ConnectionManager{
		users:  make(map[int]*UserSession, 0),
		server: server,
	}
}

func (cm *ConnectionManager) Start() {
	cm.server.Listen(cm.UserConnects)
}

func (cm *ConnectionManager) UserConnects(request *server.ConnectionRequest, addr *net.UDPAddr) {
	fmt.Printf("User Connecting: %v \n", request.UserId)
	cm.server.Write(fmt.Sprintf("You have connected! (UserId: %v)", request.UserId), addr)

	session := NewUserSession(request.UserId, &request.Friends, &Connection{Udp: addr})

	go cm.MonitorOnlineStatusRequests(session)
	go cm.MonitorOnlineStatusResponses(session)

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.users[request.UserId] = session
}

func (cm *ConnectionManager) MonitorOnlineStatusRequests(session *UserSession) {

	request := <-session.OnlineStatusRequestChan

	// if the requested id is in the ConnectionManager's list of users...
	if friendTarget, found := cm.users[request.FriendId]; found {
		// todo: concurrency
		// todo: validate symmatry

		// ...add requester to targets friends list
		friendTarget.Friends[request.Requester.UserId] = request.Requester

		// ...send the target friend back to the requester with a happy message
		session.OnlineStatusResponseChan <- &OnlineStatusResponse{
			Message: fmt.Sprintf("Friend online: %v", friendTarget.UserId),
			Friend:  friendTarget,
			Request: request,
		}
	} else {
		// ...otherwise send back the "empty value" OfflineUser and a sad message
		session.OnlineStatusResponseChan <- &OnlineStatusResponse{
			Message: "User not online",
			Friend:  OfflineUser,
			Request: request,
		}
	}
}

func (cm *ConnectionManager) MonitorOnlineStatusResponses(session *UserSession) {
	// todo: concurrency
	statusResponse := <-session.OnlineStatusResponseChan

	currentSessionFriend := session.Friends[statusResponse.Request.FriendId]
	responseFriend := statusResponse.Friend
	// if the address of these friends has changed update the requesters list and notify the updated online status
	if &currentSessionFriend != &responseFriend {
		session.Friends[statusResponse.Request.FriendId] = responseFriend
		// todo: notify
	}

}
