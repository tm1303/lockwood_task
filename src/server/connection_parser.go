package server

import (
	"encoding/json"
	"fmt"
)

type ConnectionRequest struct {
	UserId  int   `json:"user_id"`
	Friends []int `json:"friends"`
}

func NewConnectionRequest(in []byte) *ConnectionRequest {

	var jsonData ConnectionRequest
	err := json.Unmarshal(in, &jsonData)
	if err != nil {
		fmt.Println("error: ", err)
	}

	return &jsonData
}
