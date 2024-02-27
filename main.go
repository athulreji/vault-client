package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Content string `json:"content"`
	From    string `json:"From"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Get username
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1] // Trim newline

	// Connect to server
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Send username
	err = conn.WriteJSON(Message{Type: "join", From: username})
	if err != nil {
		fmt.Println("Error sending username:", err)
		return
	}

	// Listen for messages and send messages
	go readMessages(conn)
	writeMessages(conn, reader, username)
}

func readMessages(conn *websocket.Conn) {
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		var msg Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}
		fmt.Println(msg.From+":", msg.Content)
	}
}

func writeMessages(conn *websocket.Conn, reader *bufio.Reader, username string) {
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1] // Trim newline

		err := conn.WriteJSON(Message{Type: "message", From: username, Content: text})
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
