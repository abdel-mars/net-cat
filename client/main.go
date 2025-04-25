package main

import (
	"fmt"
	"net"
)
func main() {
	//connect to the server

	conn, err := net.Dial("tcp", "localhost:8989")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to the server. Sending message...")

	//send data to the server

	_, err = conn.Write([]byte("welcome\n"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Message sent. Closing connection...")

	//close the connection
	conn.Close()
}
