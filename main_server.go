package main

import "go_chat/server"

func main() {
	server := server.NewServer("127.0.0.1", 8000)
	server.Start()
}
