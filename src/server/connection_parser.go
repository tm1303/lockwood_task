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

	fmt.Printf("in:\n %v\n", string(in[:]))

	var jsonData ConnectionRequest
	err := json.Unmarshal(in, &jsonData)
	fmt.Println(err)

	return &jsonData
}
