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
	IsOnline                bool
	OnlineStatusRequestChan chan *OnlineStatusRequest
}

type OnlineStatusRequest struct {
	FriendId  int
	Requester *UserSession
}

func NewUserSession(userId int, friendIds *[]int, con *Connection, usm *UserSessionManager) *UserSession {
	friends := make(map[int]*UserSession, len(*friendIds))
	for _, friendId := range *friendIds {
		// initialy assume all friends are offline
		friends[friendId] = OfflineUser
	}

	userSession := &UserSession{
		UserId:                  userId,
		Friends:                 friends,
		Connection:              con,
		IsOnline:                true,
		OnlineStatusRequestChan: make(chan *OnlineStatusRequest),
	}

	go userSession.MonitorOnlineStatusRequests(usm)

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
		if friendTarget, found := usm.GetConnectedUser(request.FriendId); found {
			if friendTarget.ValidateFriendSymmetry(s.UserId) {
				friendTarget.UpdateFriend(request.Requester)
				s.UpdateFriend(friendTarget)
				continue
			}
		} else {
			s.FriendOffline(request.FriendId)
		}
	}
}
