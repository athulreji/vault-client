package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Message struct {
	Type       string `json:"type"`
	IsGroupMsg bool   `json:"isGroupMsg"`
	Group      string `json:"group"`
	To         string `json:"to"`
	Content    string `json:"content"`
	From       string `json:"From"`
}

var chatItems = make(map[string][]Message)

func (item Message) Title() string {
	return item.From
}

func (item Message) Description() string {
	return item.Content
}

func (item Message) FilterValue() string {
	return item.Content
}

type Chat struct {
	name, desc string
}

func (item Chat) Title() string {
	return item.name
}

func (item Chat) Description() string {
	return item.desc
}

func (item Chat) FilterValue() string {
	return item.name
}

type component int

const (
	chat component = iota
	message
	input
)

type state int

const (
	initializing state = iota
	ready
)

type model struct {
	height      int
	width       int
	chats       list.Model
	messages    list.Model
	currentChat string
	focus       component
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Message:
		if _, ok := chatItems[msg.From]; !ok {
			m.chats.InsertItem(0, Chat{msg.From, "today"})
			if m.currentChat == "" {
				m.currentChat = msg.From
			}
		}
		chatItems[msg.From] = append(chatItems[msg.From], msg)
		if m.currentChat == msg.From {
			m.messages.InsertItem(len(m.messages.Items()), msg)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "c":
			m.focus = chat
		case "m":
			m.focus = message
		case "/":
			m.focus = input
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.height = msg.Height - h
		m.width = msg.Width - v
		m.chats.SetSize(m.width/3-1, m.height-4)
		m.messages.SetSize((2*m.width/3)-1, m.height-4)
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.chats, cmd = m.chats.Update(msg)
	cmds = append(cmds, cmd)

	m.messages, cmd = m.messages.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#57288d"))

	headview := style.Width(m.width).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render(" Vault V0.1")
	inputview := style.Width(m.width).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render("----------------------------------------")
	chatsview := style.Width(m.width/3 - 1).Height(m.height - 4).Render(m.chats.View())
	messagesview := style.Width((2 * m.width / 3) - 1).Height(m.height - 4).Render(m.messages.View())
	// head := "Vault v0.1"
	return lipgloss.JoinVertical(lipgloss.Top, headview,
		lipgloss.JoinHorizontal(lipgloss.Top, chatsview, messagesview),
		inputview,
	)
}

func main() {
	m := model{chats: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0), messages: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)}
	m.chats.Title = "chats"
	m.messages.Title = "messages"

	p := tea.NewProgram(m, tea.WithAltScreen())
	go readMessagesAlt(p)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error creating bubble")
	}

}

// var chats = make(map[string][]Message)

func mainAlt() {
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
	go readMessages(conn)

	for true {
		// Listen for messages and send messages
		fmt.Println("Choose option:")
		fmt.Println("1. Direct Message")
		fmt.Println("2. Group Message")
		fmt.Println("3. Create Group")
		fmt.Println("4. Join Group")

		var option int

		fmt.Scanf("%d", &option)

		if option == 1 {
			directMessage(conn, username, reader)
		} else if option == 2 {
			groupMessage(conn, username, reader)
		} else if option == 3 {
			createGroup(conn, username, reader)
		} else {
			joinGroup(conn, username, reader)
		}
	}

}

func createGroup(conn *websocket.Conn, from string, reader *bufio.Reader) {
	fmt.Println("Enter group name:")
	group, _ := reader.ReadString('\n')
	group = group[:len(group)-1]

	writeCommands(conn, reader, from, "", group, false, "create group")
}

func joinGroup(conn *websocket.Conn, from string, reader *bufio.Reader) {
	fmt.Println("Enter group name:")
	group, _ := reader.ReadString('\n')
	group = group[:len(group)-1]

	writeCommands(conn, reader, from, "", group, false, "join group")
}

func directMessage(conn *websocket.Conn, from string, reader *bufio.Reader) {
	fmt.Println("Enter senders username:")
	to, _ := reader.ReadString('\n')
	to = to[:len(to)-1]
	writeMessages(conn, reader, from, to, "", false)
}

func groupMessage(conn *websocket.Conn, from string, reader *bufio.Reader) {
	fmt.Println("Enter group name:")
	group, _ := reader.ReadString('\n')
	group = group[:len(group)-1]
	writeMessages(conn, reader, from, "", group, true)
}
func readMessagesAlt(p *tea.Program) {
	username := "athul"

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
		// chats[msg.From] = append(chats[msg.From], msg)
		fmt.Println(msg.From+":", msg.Content)
	}
}

func writeMessages(conn *websocket.Conn, reader *bufio.Reader, from string, to string, group string, isGroupMsg bool) {
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1] // Trim newline

		err := conn.WriteJSON(Message{Type: "message", From: from, To: to, Group: group, IsGroupMsg: isGroupMsg, Content: text})
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}

func writeCommands(conn *websocket.Conn, reader *bufio.Reader, from string, to string, group string, isGroupMsg bool, content string) {
	err := conn.WriteJSON(Message{Type: "cmd", From: from, To: to, Group: group, IsGroupMsg: isGroupMsg, Content: content})
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
}
