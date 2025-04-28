# TCP-Chat Project

This is a simple TCP-based chat application written in Go. It consists of a server and a client program that communicate over TCP sockets.

## Project Structure

- `server/`
  - `main.go`: The TCP chat server implementation.
  - `welcome_logo.go`: Contains the ASCII art welcome logo displayed to clients on connection.
- `client/`
  - `main.go`: The TCP chat client implementation.

## Features

- Multiple clients can connect to the server concurrently (up to a maximum limit).
- Clients are prompted to enter their username upon connection.
- Messages sent by clients are broadcast to all other connected clients.
- Clients can quit the chat by typing `/quit`.
- The server displays a welcome logo when clients connect.

## How to Run

### Running the Server

1. Open a terminal and navigate to the project root directory.

2. Run the server by including all Go files in the `server` directory:

   ```bash
   go run server/*.go
   ```

   By default, the server listens on port `8989`.

3. To specify a custom port, pass it as an argument:

   ```bash
   go run server/*.go 12345
   ```

### Running the Client

1. Open another terminal and navigate to the project root directory.

2. Run the client with the server IP and port as arguments:

   ```bash
   go run client/main.go <server-ip> <server-port>
   ```

   For example, if the server is running locally on port 8989:

   ```bash
   go run client/main.go localhost 8989
   ```

3. When prompted, enter your username to join the chat.

4. Type messages and press Enter to send them. Type `/quit` to exit.

## Notes

- The server limits the number of concurrent clients to 10.
- The client and server communicate using simple text messages over TCP.
- The welcome logo is displayed to clients upon connection.

## Dependencies

- Go programming language (version 1.16 or higher recommended).

## License

This project is open source and free to use.
