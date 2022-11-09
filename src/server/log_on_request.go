package server

import (
	"encoding/json"
	"fmt"
)

type LogOnRequest struct {
	UserId  int   `json:"user_id"`
	Friends []int `json:"friends"`
}

func ParseLogOnRequest(in []byte) *LogOnRequest {

	var jsonData LogOnRequest
	err := json.Unmarshal(in, &jsonData)
	if err != nil {
		fmt.Println("error: ", err)
	}

	return &jsonData
}
