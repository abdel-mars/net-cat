# 🧵 TCPChat — A Go NetCat-Style Real-Time Chat App

This is a Go network programming project that recreates the core features of the classic **NetCat** tool — but with a **real-time group chat server**.

The project simulates a TCP-based client-server architecture, where one server listens and multiple clients can connect, chat, and disconnect — with specific behaviors and formatting rules.

---

## 🔍 What is This Project?

You're building your own version of **NetCat** in **Go**, but with **group chat** functionality.

It behaves like a mini real-time chat app over **TCP sockets**.

---

## 🧠 Key Concepts You Will Learn

- TCP socket programming in Go.
- Concurrency using goroutines, channels, or mutexes.
- Client-server communication via the `net` package.
- Synchronization and data sharing safely (via `sync.Mutex` or channels).
- Parsing and formatting messages.
- Handling user input/output over TCP.
- Error handling and graceful shutdowns.
- Managing multiple connections efficiently.

---

## ✅ Required Features (Detailed)

| Feature                | Explanation                                                                 |
|------------------------|-----------------------------------------------------------------------------|
| ✅ TCP Server           | Listens on a port (default: `8989`). Accepts up to 10 clients.              |
| ✅ TCP Clients          | Clients connect to the server and chat. Must provide **non-empty usernames**.|
| ✅ Message Broadcast    | Every client message is sent to **all connected clients**.                  |
| ✅ Time and User Tags   | Format: `[2025-04-24 12:00:00][Mars]: Hello world!`                         |
| ✅ Historical Messages  | When a new client joins, they receive all previous chat history.            |
| ✅ Join/Leave Notice    | Server announces when users join/leave.                                     |
| ✅ Linux-style Welcome  | Server greets new users with a Linux ASCII logo and a username prompt.      |
| ✅ Connection Limit     | Limit to **10 clients max**. Enforced by server logic.                      |
| ✅ Usage Validation     | Default port `8989`. If more than one argument: print `[USAGE]: ./TCPChat $port`|

---

## 🛠 Tools You Can Use

Only these Go packages are allowed:

io, log, os, fmt, net, sync, time, bufio, errors, strings, reflect


---

## 🧪 Bonus Features (Optional but Impressive)

- 🎨 Terminal UI using [`gocui`](https://github.com/jroimartin/gocui)
- 💾 Save chat logs to a file
- 💬 Multiple group chats by port/identifier

---

## ⏳ Time Estimation Breakdown

| Task                                             | Estimated Time        |
|--------------------------------------------------|------------------------|
| ✅ Set up server, listen on port                 | 1–2 hours              |
| ✅ Handle incoming client connections            | 2–3 hours              |
| ✅ Client name entry + validation                | 1 hour                 |
| ✅ Broadcast messages to all clients             | 2–3 hours              |
| ✅ Message formatting (timestamp + username)     | 1 hour                 |
| ✅ Message history for new clients               | 2–3 hours              |
| ✅ Join/leave notifications                      | 1 hour                 |
| ✅ Error handling + disconnects                  | 2 hours                |
| ✅ Connection limit (10 max)                     | 1 hour                 |
| ✅ Argument handling & usage message             | 30 min                 |
| 🧪 Optional: Terminal UI with gocui              | 4–6 hours              |
| 🧪 Optional: Save logs to file                   | 1–2 hours              |
| 🧪 Optional: Support multiple chat rooms         | 3–5 hours              |

**Total Time (Core Features Only):** 14–18 hours  
**With Bonuses:** 20–30 hours

---

## 🧠 Development Tips

- Use a `Client` struct to store each user's connection and metadata.
- Track active users using `map[string]*Client` or a `[]*Client`.
- Store chat history in a `[]string` or via buffered channels.
- Use `sync.Mutex` or Go channels to safely access shared resources.
- Don’t forget to flush `bufio.Writer` when writing responses.

---

## 💡 Example Architecture

```go
type Client struct {
    name string
    conn net.Conn
    msg  chan string
}

var clients map[string]*Client // Active connected users
var messages []string          // Chat history

// Use goroutines, channels, and mutexes to manage concurrent access
