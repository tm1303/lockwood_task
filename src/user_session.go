package lockwood_task

type UserId string

type UserSession struct {
	UserId
	Friends map[UserId]*UserSession
	*Connection
}
