package main

import "Pixivel/internal/server"

func main() {
	server := server.NewServer()
	server.Init()
	server.Run()
}
