package main

import "fmt"

func main() {
	server := NewServer()

	fmt.Println("Listening at port 1234")
	server.StartServer(":1234")
}
