package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"signal_client/models"
	"signal_client/services"
	"signal_client/services/newModel"
)

func main() {
	m := newModel.NewModel()
	m.Chats.Title = "Chats"
	m.Chats.Styles.TitleBar.Align(lipgloss.Center)

	p := tea.NewProgram(m, tea.WithAltScreen())
	channel := make(chan *websocket.Conn)
	go services.StartConnection(p, channel)
	models.ServerConn = <-channel
	if _, err := p.Run(); err != nil {
		fmt.Println("Error creating bubble")
	}
}
