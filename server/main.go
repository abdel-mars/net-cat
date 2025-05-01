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

		// Check if we have reached the maximum number of clients.
		mutex.Lock()
		if len(clients) >= maxClients {
			conn.Write([]byte("Server full. Try again later.\n"))
			conn.Close()
			mutex.Unlock()
			continue
		}
		mutex.Unlock()

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

	// Create a new Client instance.
	client := &Client{conn: conn, name: name}

	// Add the client to the clients map.
	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	// Announce the new client to all connected clients.
	broadcast(fmt.Sprintf("%s has joined our chat...", client.name), "")

	// Send the chat history to the newly connected client.
	sendHistory(conn)

	// Continuously read messages from the client.
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue // Ignore empty messages.
		}
		// Broadcast the received message to all clients.
		broadcast(text, client.name)
	}

	// Remove the client from the clients map on disconnection.
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()

	// Announce the client's departure to all connected clients.
	broadcast(fmt.Sprintf("%s has left our chat...", client.name), "")
}

// readName reads a client's name from the provided connection.
// It continuously prompts the client until a non-empty name is provided.
func readName(conn net.Conn) (string, error) {
	// Create a buffered reader from the connection.
	reader := bufio.NewReader(conn)

	// Continuously read the client's name until a valid name is provided.
	for {
		// Read the client's name.
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// Trim any whitespace characters from the name.
		name = strings.TrimSpace(name)

		// If the name is not empty, return it.
		if name != "" {
			return name, nil
		}

		// Prompt the client to enter their name again since the name is empty.
		conn.Write([]byte("[ENTER YOUR NAME]:\n"))
	}
}

// broadcast sends the provided message to all connected clients.
// If the senderName is empty, it sends the message as a server message.
// Otherwise, it formats the message with the sender's name and current time.
func broadcast(message, senderName string) {
	mutex.Lock()
	defer mutex.Unlock()

	var formatted string
	if senderName == "" {
		// If the senderName is empty, the message is sent as a server message.
		formatted = message
	} else {
		// Format the message with the sender's name and current time.
		formatted = fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), senderName, message)
		// Store the formatted message in the messages slice for later use.
		messages = append(messages, formatted)
	}

	// Send the formatted message to all connected clients.
	for _, client := range clients {
		fmt.Fprintln(client.conn, formatted)
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
