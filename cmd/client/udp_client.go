package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	lockwood_server "lockwood_task/src/server"
)

func main() {

	payload := gatherPayload()

	s, err := net.ResolveUDPAddr("udp4", "localhost:13131")
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		_, err = c.Write(payload)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("now sleep: %v \n", time.Now().Second())
		time.Sleep(3 * time.Second)

	}
}

// todo: validate these inputs, but not today.
// todo: handle some err
// todo: extract and share with tcp client (todo!)
func gatherPayload() []byte {

	fmt.Print("\n\n\nPlease enter the following carefully\n")
	fmt.Print("Validation not complete!\n\n")
	reader := bufio.NewReader(os.Stdin)

	var userId int
	for {
		fmt.Print("Connect as user (int): ")
		userId_in, _ := reader.ReadString('\n')
		userIdCount, _ := fmt.Sscan(userId_in, &userId)
		if userIdCount == 1 {
			break
		}
	}

	friends := make([]int, 0)
	fmt.Print("Enter friend IDs one at a time, leave empty for no more friends\n")
	var friendId int
	for {
		fmt.Print("   ...Friend Id: (int, or empty): ")
		friendId_in, _ := reader.ReadString('\n')
		intCount, _ := fmt.Sscan(friendId_in, &friendId)
		if intCount == 0 {
			break
		}
		friends = append(friends, friendId)
	}

	request := lockwood_server.ConnectionRequest{UserId: userId, Friends: friends} // don't like this coupling but time is short
	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("error:", err)
	}

	return payload
}
