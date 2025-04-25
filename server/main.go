package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type client struct {
	conn     net.Conn
	username string
}

type message struct {
	text string
	time time.Time
	user string
}

type server struct {
	clients     []*client
	messages    []message
	mu          sync.Mutex
}

func getWelcomeLogo() string {
	return `
	Welcome to TCP-Chat!
	         _nnnn_
	        dGGGGMMb
	       @p~qp~~qMb
	       M|@||@) M|
	       @,----.JM|
	      JS^\__/  qKL
	     dZP        qKRb
	    dZP          qKKb
	   fZP            SMMb
	   HZM            MMMM
	   FqM            MMMM
	 __| ".        |\dS"qML
	 |    '.       | '' \Zq
	_)      \.___.,|     .'
	\____   )MMMMMP|   .'
	     '-'       '--'
	`
}

func (s *server) handleClient(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("\n[ENTER YOUR NAME]: "))

	scanner := bufio.NewScanner(conn)
	var username string
	for scanner.Scan() {
		username = scanner.Text()
		if username != "" {
			break
		}
		conn.Write([]byte("[ENTER YOUR NAME]: "))
	}
	if username == "" {
		return
	}

	newClient := &client{conn: conn, username: username}
	s.mu.Lock()
	if len(s.clients) >= 10 {
		conn.Write([]byte("Chat is full. Disconnecting.\n"))
		s.mu.Unlock()
		return
	}
	s.clients = append(s.clients, newClient)
	s.mu.Unlock()

	s.mu.Lock()
	for _, msg := range s.messages {
		conn.Write([]byte(msg.text + "\n"))
	}
	s.mu.Unlock()

	s.broadcast(fmt.Sprintf("[%s][Server]: %s has joined the chat!\n", 
		time.Now().Format("2006-01-02 15:04:05"), username))

	input := bufio.NewScanner(conn)
	for input.Scan() {
		text := input.Text()
		if text == "" {
			continue
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		msg := fmt.Sprintf("[%s][%s]: %s", timestamp, username, text)

		s.mu.Lock()
		s.messages = append(s.messages, message{text: msg, time: time.Now(), user: username})
		s.mu.Unlock()
		s.broadcast(msg + "\n")
	}

	s.mu.Lock()
	for i, c := range s.clients {
		if c == newClient {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
	s.broadcast(fmt.Sprintf("[%s][Server]: %s has left the chat.\n", 
		time.Now().Format("2006-01-02 15:04:05"), username))
}

func (s *server) broadcast(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range s.clients {
		_, err := c.conn.Write([]byte(msg))
		if err != nil {
			log.Printf("Error broadcasting to %s: %v\n", c.username, err)
		}
	}
}

func main() {
	port := "8989"
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	s := &server{
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Listening on port :%s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handleClient(conn)
	}
}
