package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// func main() {
// 	//listen on port

// 	ln, err := net.Listen("tcp", ":8989")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		//handle connection theads
// 		go HandleConn(conn)
// 	}
// }
// // Handles a single connection.
// func HandleConn(conn net.Conn) {
// 	//Read the data from the socket
// 	for {
// 		data, err := bufio.NewReader(conn).ReadString('\n')
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		fmt.Println(data)
// 	}
// 	//close the connection
// 	conn.Close()
// }

type client struct {
	name string
	conn net.Conn
}

type message struct {
	user string
	text string
	time time.Time
}

type Server struct {
	clients []*client
	messages []message
	mu sync.Mutex
}

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("[ENTER YOUR NAME]: "))
	scanner := bufio.NewScanner(conn)
	var username string
	for scanner.Scan() {
		username = scanner.Text()
		if username != "" {
			break
		}
		conn.Write([]byte("[ENTER YOUR NAME]: "))
	}
	//logout user if didnt put a name

	if username == "" {
		return
	}

	newClient := &client {conn : conn, name : username}
	s.mu.Lock()
	if len(s.clients) > 10 {
		conn.Write([]byte ("Chat is full. Disconnecting.\n"))
		s.mu.Unlock()
		return
	}
	s.clients = append(s.clients, newClient)
	s.mu.Unlock()



}

func main() {
	port := "8989"
	if len(os.Args) > 2 {
		fmt.Println("Usage: ./tcpchat <port>")
	}
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	s := &Server{}

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.HandleConn(conn)
	}
}