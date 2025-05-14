package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultPort = "8989"
const maxClients = 10

type Client struct {
	conn net.Conn
	name string
}

var (
	clients  = make(map[net.Conn]*Client)
	messages []string
	mutex    sync.Mutex
)

// main runs the TCP-Chat server.
func main() {
	// Get the port to listen on.
	port, err := getPort()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen on the specified port.
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	defer listener.Close()

	// Print the port number.
	fmt.Println("Listening on port:", port)

	// Run an infinite loop to handle connections.
	for {
		// Accept a connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept:", err)
			continue
		}

		// Handle the connection in a separate goroutine.
		go handleConnection(conn)
	}
}

// getPort retrieves the port number from the command-line arguments.
// If no port is provided, it returns the default port.
// If an invalid port is provided, it returns an error.
func getPort() (string, error) {
	args := os.Args[1:]

	// No port provided, return the default port.
	if len(args) == 0 {
		return defaultPort, nil
	}	
	// One argument provided, validate if it's a valid integer port.
	if len(args) == 1 {
			_, err := strconv.Atoi(args[0])
			if err != nil {
				return "", fmt.Errorf("[USAGE]: ./server $port")
			}
			return args[0], nil
		}
	
		// Invalid number of arguments, return an error.
		return "", fmt.Errorf("[USAGE]: ./server $port")
	}
	
// handleConnection manages a single client connection to the chat server.
func handleConnection(conn net.Conn) {
	// Ensure the connection is closed when the function exits.
	defer conn.Close()

	// Send the welcome logo to the client.
	conn.Write([]byte(welcomeLogo))

	// Read the client's name.
	name, err := readName(conn)
	if err != nil {
		return
	}

	// Added maxClients check here after username is provided
	mutex.Lock()
	if len(clients) >= maxClients {
		conn.Write([]byte("Server full. Try again later.\n"))
		mutex.Unlock()
		return
	}
	mutex.Unlock()

	// Create a new Client instance.
	client := &Client{conn: conn, name: name}

	// Add the client to the clients map.
	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	// Announce the new client to all connected clients.
	broadcast(fmt.Sprintf("\033[32m%s has joined our chat...\033[0m\n", client.name), "")

	// Send the chat history to the newly connected client.
	sendHistory(conn)

	// Continuously read messages from the client.
	scanner := bufio.NewScanner(conn)
	// Write the prefix before reading the first message
	conn.Write([]byte(Prefix(client.name)))
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			// Write the prefix again for empty input
			conn.Write([]byte(Prefix(client.name)))
			continue // Ignore empty messages.
		}

		if !isPrintableASCII(text) {
			conn.Write([]byte("[MESSAGE MUST CONTAIN ONLY PRINTABLE ASCII CHARACTERS]\n"))
			// Write the prefix again after error message
			conn.Write([]byte(Prefix(client.name)))
			continue
		}

		// Broadcast the received message to all clients.
		broadcast(text, client.name)
		// Write the prefix again after broadcasting
		conn.Write([]byte(Prefix(client.name)))
	}

	// Remove the client from the clients map on disconnection.
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()

	// Announce the client's departure to all connected clients.
	broadcast(fmt.Sprintf("\033[31m%s has left our chat...\033[0m\n", client.name), "")
}

func isPrintableASCII(s string) bool {
	for _, r := range s {
		if r < 32 || r > 127 {
			return false
		}
	}
	return true
}

func readName(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		name = strings.TrimSpace(name)

		if name == "" {
			conn.Write([]byte("[ENTER YOUR NAME]:"))
			continue
		}

		if !isPrintableASCII(name) {
			conn.Write([]byte("[NAME MUST CONTAIN ONLY PRINTABLE ASCII CHARACTERS] \n[ENTER YOUR NAME]: "))
			continue
		}

		// Check if the name is already taken
		mutex.Lock()
		nameTaken := false
		for _, client := range clients {
			if client.name == name {
				nameTaken = true
				break
			}
		}
		mutex.Unlock()

		if nameTaken {
			conn.Write([]byte("[NAME ALREADY TAKEN] \n[ENTER YOUR NAME]: "))
			continue
		}

		return name, nil
	}
}

// broadcast sends the provided message to all connected clients.
// If the senderName is empty, it sends the message as a server message.
// Otherwise, it formats the message with the sender's name and current time.
// broadcast sends the provided message to all connected clients.
func Prefix(name string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:", timestamp, name)
}

func broadcast(message, senderName string) {
	mutex.Lock()
	defer mutex.Unlock()

	var formatted string
	if senderName == "" {
		formatted = message
	} else {
		formatted = Prefix(senderName) + message
		messages = append(messages, formatted)
	}

	// Find sender's connection
	var senderConn net.Conn
	for conn, client := range clients {
		if client.name == senderName {
			senderConn = conn
			break
		}
	}

	// Send message to all clients except the sender
	for conn := range clients {
		if senderConn != nil && conn == senderConn {
			continue
		}
		// Write escape sequences to save cursor, clear line, and restore cursor to fix input line display
		conn.Write([]byte("\033[s\033[2K\r")) // Clear line and move cursor.
		fmt.Fprintln(conn, formatted)
		// Write the prefix after the message for all clients
		conn.Write([]byte(Prefix(clients[conn].name)))
		conn.Write([]byte("\033[u\033[B")) // Restore cursor and move down.
	}
}
// sendHistory sends the chat history to a newly connected client.
func sendHistory(conn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	// Iterate through all the stored messages and send them to the newly connected client.
	for _, msg := range messages {
		fmt.Fprintln(conn, msg)
	}
}
