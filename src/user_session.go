package lockwood_task

var OfflineUser *UserSession = &UserSession{
	UserId:     -1,
	Friends:    nil,
	Connection: nil,
	IsOnline:   false,
	// todo : default chans to handle always rejected requests
}

type UserSession struct {
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

func NewUserSession(userId int, friendIds *[]int, con *Connection) *UserSession {
	friends := make(map[int]*UserSession, len(*friendIds))
	for _, friendId := range *friendIds {
		// assume this user is not online
		friends[friendId] = OfflineUser
	}

	return &UserSession{
		UserId:                   userId,
		Friends:                  friends,
		Connection:               con,
		IsOnline:                 true,
		OnlineStatusRequestChan:  make(chan *OnlineStatusRequest),
		OnlineStatusResponseChan: make(chan *OnlineStatusResponse),
	}
}
