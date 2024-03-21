package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

var (
	chatItems  = make(map[string][]Message)
	username   = ""
	serverConn *websocket.Conn
)

var (
	titleStyle        = lipgloss.NewStyle().PaddingLeft(0).PaddingTop(0).Foreground(lipgloss.Color("#bdae93")).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#689d6a"))
	noItemStyle       = lipgloss.NewStyle().PaddingLeft(1).PaddingTop(2).Foreground(lipgloss.Color("#689d6a")).Bold(false)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#fabd2f"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

func main() {
	chatlist := list.New([]list.Item{}, chatItemDelegate{}, 0, 0)
	chatlist.SetShowHelp(false)
	chatlist.Styles = list.Styles{
		TitleBar:        titleStyle,
		NoItems:         noItemStyle,
		PaginationStyle: paginationStyle,
	}
	chatlist.SetShowFilter(false)
	chatlist.SetShowStatusBar(false)

	input := textinput.New()
	usernameInput := textinput.New()
	passwordInput := textinput.New()
	groupnameInput := textinput.New()

	messagelist := list.New([]list.Item{}, messageItemDelegate{}, 0, 0)
	messagelist.Styles = list.Styles{
		NoItems: lipgloss.NewStyle().PaddingLeft(2).PaddingTop(2),
	}
	messagelist.SetShowHelp(false)
	messagelist.SetShowTitle(false)
	messagelist.SetShowStatusBar(false)
	messagelist.SetShowFilter(false)
	m := model{
		chats:          chatlist,
		messages:       messagelist,
		input:          input,
		usernameInput:  usernameInput,
		passwordInput:  passwordInput,
		groupnameInput: groupnameInput,
	}
	m.chats.Title = "Chats"
	m.chats.Styles.TitleBar.Align(lipgloss.Center)

	p := tea.NewProgram(m, tea.WithAltScreen())
	channel := make(chan *websocket.Conn)
	go startConnection(p, channel)
	serverConn = <-channel
	if _, err := p.Run(); err != nil {
		fmt.Println("Error creating bubble")
	}
}

func userLogin() {
	// Send username
	err := serverConn.WriteJSON(Message{Type: "cmd", From: username, Content: "login"})
	if err != nil {
		fmt.Println("Error sending username:", err)
		return
	}
}

func startConnection(p *tea.Program, ch chan *websocket.Conn) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	ch <- conn
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
		p.Send(msg)
		// chats[msg.From] = append(chats[msg.From], msg)
		// fmt.Println(msg.From+":", msg.Content)
	}
}

func writeMessages(conn *websocket.Conn, msg Message) {
	err := conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
}

// func createGroup(conn *websocket.Conn, msg Message) {
// 	// Send username
// 	err := serverConn.WriteJSON(msg)
// 	if err != nil {
// 		fmt.Println("Error sending username:", err)
// 		return
// 	}

// 	writeCommands(conn, reader, from, "", group, false, "create group")
// }

// func joinGroup(conn *websocket.Conn, from string, reader *bufio.Reader) {
// 	fmt.Println("Enter group name:")
// 	group, _ := reader.ReadString('\n')
// 	group = group[:len(group)-1]

// 	writeCommands(conn, reader, from, "", group, false, "join group")
// }

// func directMessage(conn *websocket.Conn, from string, reader *bufio.Reader) {
// 	fmt.Println("Enter senders username:")
// 	to, _ := reader.ReadString('\n')
// 	to = to[:len(to)-1]
// 	writeMessages(conn, reader, from, to, "", false)
// }

// func groupMessage(conn *websocket.Conn, from string, reader *bufio.Reader) {
// 	fmt.Println("Enter group name:")
// 	group, _ := reader.ReadString('\n')
// 	group = group[:len(group)-1]
// 	writeMessages(conn, reader, from, "", group, true)
// }

// func writeCommands(conn *websocket.Conn, reader *bufio.Reader, from string, to string, group string, isGroupMsg bool, content string) {
// 	err := conn.WriteJSON(Message{Type: "cmd", From: from, To: to, Group: group, IsGroupMsg: isGroupMsg, Content: content})
// 	if err != nil {
// 		fmt.Println("Write error:", err)
// 		return
// 	}
// }

// func readMessages(conn *websocket.Conn) {
// 	for {
// 		_, msgBytes, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Read error:", err)
// 			return
// 		}

// 		var msg Message
// 		err = json.Unmarshal(msgBytes, &msg)
// 		if err != nil {
// 			fmt.Println("Error decoding JSON:", err)
// 			continue
// 		}
// 		// chats[msg.From] = append(chats[msg.From], msg)
// 		fmt.Println(msg.From+":", msg.Content)
// 	}
// }
