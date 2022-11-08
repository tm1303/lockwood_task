package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	s, err := net.ResolveUDPAddr("udp4", "localhost:13131")
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		tx := fmt.Sprintf(`{
				"user_id": %v, 
				"friends": [2, 3, 4]
				}`, time.Now().Second())

		data := []byte(tx)
		_, err = c.Write(data)

		if err != nil {
			fmt.Println(err)
			//return
		}

		fmt.Printf("now sleep: %v \n", time.Now().Second())
		time.Sleep(3 * time.Second)

	}
}
