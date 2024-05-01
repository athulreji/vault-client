package services

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"net/url"
	"signal_client/models"
	"signal_client/models/message"
)

func StartConnection(p *tea.Program, ch chan *websocket.Conn) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	ch <- conn
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		var msg message.Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}
		p.Send(msg)
	}
}

func WriteMessages(conn *websocket.Conn, msg message.Message) {
	err := conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
}

func UserLogin() {
	// Send username
	err := models.ServerConn.WriteJSON(message.Message{Type: "cmd", From: models.Username, Content: "login"})
	if err != nil {
		fmt.Println("Error sending username:", err)
		return
	}
}
