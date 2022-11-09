package lockwood_task

import (
	// "fmt"
	"fmt"
	"lockwood_task/src/server"
	"sync"
	"time"
)

var OfflineUser *UserSession = &UserSession{
	UserId:   -1,
	Friends:  nil,
	Notifier: nil,
	IsOnline: false,
	// todo : default chans to handle always rejected requests
}

type UserSession struct {
	UserId                  int
	Notifier                server.UserNotifierChannel
	IsOnline                bool
	OnlineStatusRequestChan chan *OnlineStatusRequest
	mutex                   sync.Mutex
	Friends                 map[int]*UserSession
}

type OnlineStatusRequest struct {
	FriendId  int
	Requester *UserSession
}

var refreshDelay time.Duration = 10 * time.Second

func NewUserSession(userId int, friendIds *[]int, notifier server.UserNotifierChannel, usm *UserSessionManager) *UserSession {
	friends := make(map[int]*UserSession, len(*friendIds))
	for _, friendId := range *friendIds {
		// initialy assume all friends are offline
		friends[friendId] = OfflineUser
	}

	userSession := &UserSession{
		UserId:                  userId,
		Friends:                 friends,
		Notifier:                notifier,
		IsOnline:                true,
		OnlineStatusRequestChan: make(chan *OnlineStatusRequest),
	}

	go userSession.MonitorOnlineStatusRequests(usm)
	go func() {
		for {
			userSession.RefreshAllFriendsOnlineStatus()
			time.Sleep(refreshDelay)
		}
	}()

	return userSession
}

func (s *UserSession) RefreshAllFriendsOnlineStatus() {
	fmt.Printf("\nRefresh all friends for %v\n", s.UserId)
	for friendId := range s.Friends {
		s.OnlineStatusRequestChan <- &OnlineStatusRequest{
			FriendId:  friendId,
			Requester: s,
		}
	}
}

func (s *UserSession) MonitorOnlineStatusRequests(usm *UserSessionManager) {
	for {
		request := <-s.OnlineStatusRequestChan
		if friendTarget, found := usm.GetConnectedUser(request.FriendId); found {
			if friendTarget.ValidateFriendRequestSymmetry(s.UserId) {
				fmt.Printf("%v accepeted %v\n", request.FriendId, s.UserId)
				friendTarget.UpdateFriend(request.Requester)
				s.UpdateFriend(friendTarget)
				continue
			}

			fmt.Printf("%v did NOT accepet %v\n", request.FriendId, s.UserId)
		} else {
			fmt.Printf("%v appears OFFLINE to %v\n", request.FriendId, s.UserId)
			s.SetFriendAsOffline(request.FriendId)
		}
	}
}

func (s *UserSession) ValidateFriendRequestSymmetry(requesterId int) bool {
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
		s.Notifier <- fmt.Sprintf("Your friend is ONLINE! (UserId: %v)", updateFriend.UserId)
	}
}

func (s *UserSession) SetFriendAsOffline(friendId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentSessionFriend := s.Friends[friendId]
	// if the address of these friends has changed update the requesters list and notify the updated online status
	if currentSessionFriend != OfflineUser {
		s.Friends[friendId] = OfflineUser
		s.Notifier <- fmt.Sprintf("Your friend is OFFLINE! (UserId: %v)", friendId)
	}
}
