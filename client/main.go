package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func readServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Disconnected from server.")
			return
		}
		fmt.Print("\r" + msg + "> ")
	}
}

func getUsername(scanner *bufio.Scanner, conn net.Conn) string {
	for {
		fmt.Print("Enter your username: ")
		scanner.Scan()
		username := strings.TrimSpace(scanner.Text())
		if username == "" {
			continue
		}
		_, err := fmt.Fprintf(conn, "%s\n", username)
		if err != nil {
			log.Fatal("Error sending username:", err)
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal("Error reading server response:", err)
		}
		if strings.Contains(response, "joined") {
			return username
		}
		fmt.Print(response)
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./client <server-address>")
		return
	}

	serverAddr := os.Args[1]
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	go readServerMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	username := getUsername(scanner, conn)

	fmt.Printf("\n[%s][System]: Welcome, %s! Start chatting.\n> ", 
		time.Now().Format("2006-01-02 15:04:05"), username)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			fmt.Print("> ")
			continue
		}
		_, err := fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
		fmt.Print("> ")
	}
}