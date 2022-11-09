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

	go func(){
		// todo: move to user_session with some funky func
		statusRequest := <- session.OnlineStatusRequestChan
		statusResponse := cm.RequestOnlineStatus(statusRequest)
		session.OnlineStatusResponseChan  <- statusResponse
	}()

	go func(){
		// todo: move to user_session with some funky func
		statusResponse := <- session.OnlineStatusResponseChan
		cm.ApplyStatusResponse(statusResponse)
		session.OnlineStatusResponseChan  <- statusResponse
	}()
	
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.users[request.UserId] = session
}

func (cm *ConnectionManager) RequestOnlineStatus(request *OnlineStatusRequest) *OnlineStatusResponse {

	if friendTarget, found := cm.users[request.FriendId]; found {

		// todo: concurrency
		// todo: validate symmatry

		// add requester to targets friends list
		friendTarget.Friends[request.Requester.UserId] = request.Requester

		return &OnlineStatusResponse{
			Message:  fmt.Sprintf("Friend online: %v", friendTarget.UserId),
			Friend: friendTarget,
			Requester: request.Requester, // todo: can probs remove requester and just grab it in scope elsewhere
		}
	} else{
		return &OnlineStatusResponse{
			Message:   "User not online",
			Friend: OfflineUser,
			Requester: request.Requester,
		}
	}
}

func (cm *ConnectionManager) ApplyStatusResponse(response *OnlineStatusResponse) {
	response.Requester.Friends[response.Friend.UserId] = response.Friend
}
