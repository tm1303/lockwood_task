package lockwood_task

import (
	"fmt"
	"sync"
)

var OfflineUser *UserSession = &UserSession{
	UserId:     -1,
	Friends:    nil,
	Connection: nil,
	IsOnline:   false,
	// todo : default chans to handle always rejected requests
}

type UserSession struct {
	mutex   *sync.Mutex
	UserId  int
	Friends map[int]*UserSession
	*Connection
	IsOnline                 bool
	OnlineStatusRequestChan  chan *OnlineStatusRequest
	OnlineStatusResponseChan chan *OnlineStatusResponse
}

type OnlineStatusRequest struct {
	FriendId  int
	Requester *UserSession
}

type OnlineStatusResponse struct {
	Message string
	Friend  *UserSession
	Request *OnlineStatusRequest
}

func NewUserSession(userId int, friendIds *[]int, con *Connection, usm *UserSessionManager) *UserSession {
	friends := make(map[int]*UserSession, len(*friendIds))
	for _, friendId := range *friendIds {
		// assume this user is not online
		friends[friendId] = OfflineUser
	}

	userSession := &UserSession{
		UserId:                   userId,
		Friends:                  friends,
		Connection:               con,
		IsOnline:                 true,
		OnlineStatusRequestChan:  make(chan *OnlineStatusRequest),
		OnlineStatusResponseChan: make(chan *OnlineStatusResponse),
	}

	go userSession.MonitorOnlineStatusRequests(usm)
	go userSession.MonitorOnlineStatusResponses(usm)

	return userSession
}

func (s *UserSession) MonitorOnlineStatusRequests(usm *UserSessionManager) {

	for {
		request := <-s.OnlineStatusRequestChan

		// if the requested id is in the ConnectionManager's list of users...
		if friendTarget, accepted := usm.GetConnectedUser(request.FriendId); accepted {
			// todo: concurrency
			// todo: validate symmatry

			// ...add requester to targets friends list
			// todo: encapsualte
			friendTarget.Friends[request.Requester.UserId] = request.Requester

			// ...send the target friend back to the requester with a happy message
			s.OnlineStatusResponseChan <- &OnlineStatusResponse{
				Message: fmt.Sprintf("Friend online: %v", friendTarget.UserId),
				Friend:  friendTarget,
				Request: request,
			}
		} else {
			// ...otherwise send back the "empty value" OfflineUser and a sad message
			s.OnlineStatusResponseChan <- &OnlineStatusResponse{
				Message: "User not online",
				Friend:  OfflineUser,
				Request: request,
			}
		}
	}
}

func (s *UserSession) MonitorOnlineStatusResponses(usm *UserSessionManager) {
	for {
		// todo: concurrency
		statusResponse := <-s.OnlineStatusResponseChan

		currentSessionFriend := s.Friends[statusResponse.Request.FriendId]
		responseFriend := statusResponse.Friend
		// if the address of these friends has changed update the requesters list and notify the updated online status
		if &currentSessionFriend != &responseFriend {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			s.Friends[statusResponse.Request.FriendId] = responseFriend
			// todo: notify
		}
	}
}
