package lockwood_task

import (
	// "fmt"
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
	// OnlineStatusResponseChan chan *OnlineStatusResponse
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
		// assume all friends are offline
		friends[friendId] = OfflineUser
	}

	userSession := &UserSession{
		UserId:                  userId,
		Friends:                 friends,
		Connection:              con,
		IsOnline:                true,
		OnlineStatusRequestChan: make(chan *OnlineStatusRequest),
		////
		//// todo: I think we can kill this second chan, bonus!
		////
		// OnlineStatusResponseChan: make(chan *OnlineStatusResponse),
	}

	go userSession.MonitorOnlineStatusRequests(usm)
	////
	//// todo: I think we can kill this second chan, bonus!
	////
	// go userSession.MonitorOnlineStatusResponses(usm)

	return userSession
}

func (s *UserSession) ValidateFriendSymmetry(requesterId int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	_, ok := s.Friends[requesterId]
	return ok
}

func (s *UserSession) UpdateFriend(updateFriend *UserSession) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentSessionFriend := s.Friends[updateFriend.UserId]
	// if the address of these friends has changed update the requesters list and notify the updated online status
	if currentSessionFriend != updateFriend {
		s.Friends[updateFriend.UserId] = updateFriend
		// todo: notify
	}
}

func (s *UserSession) FriendOffline(friendId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentSessionFriend := s.Friends[friendId]
	// if the address of these friends has changed update the requesters list and notify the updated online status
	if currentSessionFriend != OfflineUser {
		s.Friends[friendId] = OfflineUser
		// todo: notify
	}
}

func (s *UserSession) MonitorOnlineStatusRequests(usm *UserSessionManager) {

	for {
		request := <-s.OnlineStatusRequestChan

		// if the requested friend id is in the ConnectionManager's list of users...
		if friendTarget, found := usm.GetConnectedUser(request.FriendId); found {
			// ...and the requested friend accepts the requester
			if friendTarget.ValidateFriendSymmetry(s.UserId) {

				// todo: concurrency
				// todo: encapsualte
				// ...add requester to targets friends list
				friendTarget.UpdateFriend(request.Requester)
				s.UpdateFriend(friendTarget)

				////
				//// todo: I think we can kill this second chan, bonus!
				////
				// ...send the target friend back to the requester with a happy message
				// s.OnlineStatusResponseChan <- &OnlineStatusResponse{
				// 	Message: fmt.Sprintf("Friend online: %v", friendTarget.UserId),
				// 	Friend:  friendTarget,
				// 	Request: request,
				// }

				continue
			}
		}

		s.FriendOffline(request.FriendId)
		// ...otherwise send back the "empty value" OfflineUser and a sad message
		// s.OnlineStatusResponseChan <- &OnlineStatusResponse{
		// 	Message: "User not online",
		// 	Friend:  OfflineUser,
		// 	Request: request,
		// }
	}
}

// //
// // todo: I think we can kill this second chan, bonus!
// //
// todo: better name please!
// func (s *UserSession) MonitorOnlineStatusResponses(usm *UserSessionManager) {
// 	for {
// 		// todo: concurrency
// 		statusResponse := <-s.OnlineStatusResponseChan
// 		responseFriend := statusResponse.Friend

// 	}
// }
