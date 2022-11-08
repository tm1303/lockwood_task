package main

import (
	// i don't know the best way to do go imports/packages/namespaces :/
	lockwood_manager "lockwood_task/src"
	lockwood_server "lockwood_task/src/server"
)

var udpPort string = ":13131"

func main() {
	server := lockwood_server.NewUdpServer(udpPort)
	connectionManager := lockwood_manager.NewConnectionManager(server)
	connectionManager.Start()
}
