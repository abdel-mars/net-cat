package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Ensure the user has provided both IP and port as command-line arguments
	// Check if the user passed the IP and port via command-line arguments
	if len(os.Args) != 3 {
		fmt.Println("[USAGE]: ./client $IP $PORT")
		return
	}

	ip := os.Args[1]
	port := os.Args[2]
	serverAddr := fmt.Sprintf("%s:%s", ip, port)

	// Connect to the server using the provided IP and port
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Create buffered readers for server and standard input
	serverReader := bufio.NewReader(conn)
	stdinReader := bufio.NewReader(os.Stdin)

	// Read and print the welcome message from the server
	welcomeMsg, err := serverReader.ReadString(':')
	if err != nil {
		fmt.Println("Error reading welcome:", err)
		return
	}
	fmt.Print(welcomeMsg)

	// Loop to read and send the user's name until accepted by the server
	for {
		name, _ := stdinReader.ReadString('\n')
		name = strings.TrimSpace(name)

		if name == "" {
			fmt.Println("[NAME CANNOT BE EMPTY] \n [ENTER YOUR NAME]:")
			continue
		}

		_, err = conn.Write([]byte(name + "\n"))
		if err != nil {
			fmt.Println("Error sending name:", err)
			return
		}

		// Read server response after sending name
		response, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server.")
			return
		}

		// If server says name is taken, prompt again
		if strings.Contains(response, "[NAME ALREADY TAKEN]") {
			fmt.Print(response)
			continue
		}

		// Otherwise, print the server response and break the loop
		fmt.Print(response)
		break
	}

	// Goroutine to constantly listen for messages from the server
	go func() {
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Disconnected from server.")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	// Main loop to read user input and send messages to the server
	for {
		// fmt.Print("You: ")
		msg, _ := stdinReader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		// Skip sending empty messages
		if msg == "" {
			continue
		}

		// Send the user's message to the server
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}
