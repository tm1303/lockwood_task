package lockwood_task

import (
	"fmt"
	"lockwood_task/src/server"
	"sync"
	"time"
)

var offlineUser *UserSession = &UserSession{
	userId:   -1,
	friends:  nil,
	Notifier: nil,
	isOnline: false,
	// todo : default chans to handle always rejected requests
}

type UserSession struct {
	userId                  int
	Notifier                server.UserNotifierChannel
	isOnline                bool
	onlineStatusRequestChan chan *OnlineStatusRequest
	mutex                   sync.Mutex
	friends                 map[int]*UserSession
}

type OnlineStatusRequest struct {
	friendId  int
	requester *UserSession
}

var refreshDelay time.Duration = 10 * time.Second

func NewUserSession(userId int, friendIds *[]int, notifier server.UserNotifierChannel, usm *UserSessionManager) *UserSession {
	friends := make(map[int]*UserSession, len(*friendIds))
	for _, friendId := range *friendIds {
		// initialy assume all friends are offline
		friends[friendId] = offlineUser
	}

	userSession := &UserSession{
		userId:                  userId,
		friends:                 friends,
		Notifier:                notifier,
		isOnline:                true,
		onlineStatusRequestChan: make(chan *OnlineStatusRequest),
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
	fmt.Printf("\nRefresh all friends for %v\n", s.userId)
	for friendId := range s.friends {
		s.onlineStatusRequestChan <- &OnlineStatusRequest{
			friendId:  friendId,
			requester: s,
		}
	}
}

func (s *UserSession) MonitorOnlineStatusRequests(usm *UserSessionManager) {
	for {
		request := <-s.onlineStatusRequestChan
		if friendTarget, found := usm.GetConnectedUser(request.friendId); found {
			if friendTarget.ValidateFriendRequestSymmetry(s.userId) {
				fmt.Printf("%v accepeted %v\n", request.friendId, s.userId)
				friendTarget.UpdateFriend(request.requester)
				s.UpdateFriend(friendTarget)
				continue
			}

			fmt.Printf("%v did NOT accepet %v\n", request.friendId, s.userId)
		} else {
			fmt.Printf("%v appears OFFLINE to %v\n", request.friendId, s.userId)
			s.SetFriendAsOffline(request.friendId)
		}
	}
}

func (s *UserSession) ValidateFriendRequestSymmetry(requesterId int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.friends[requesterId]
	return ok
}

func (s *UserSession) UpdateFriend(updateFriend *UserSession) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentSessionFriend := s.friends[updateFriend.userId]
	// if the address of these friends has changed update the requesters list and notify the updated online status
	if currentSessionFriend != updateFriend {
		s.friends[updateFriend.userId] = updateFriend
		s.Notifier <- fmt.Sprintf("Your friend is ONLINE! (UserId: %v)", updateFriend.userId)
	}
}

func (s *UserSession) SetFriendAsOffline(friendId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentSessionFriend := s.friends[friendId]
	// if the address of these friends has changed update the requesters list and notify the updated online status
	if currentSessionFriend != offlineUser {
		s.friends[friendId] = offlineUser
		s.Notifier <- fmt.Sprintf("Your friend is OFFLINE! (UserId: %v)", friendId)
	}
}
