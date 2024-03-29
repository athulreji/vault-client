package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type quit bool

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	// Process Incoming Message
	case Message:
		recvChat := Chat{}
		if msg.IsGroupMsg {
			recvChat.isGroupChat = true
			recvChat.name = msg.Group
		} else {
			recvChat.isGroupChat = false
			recvChat.name = msg.From
		}
		if _, ok := chatItems[recvChat.name]; !ok {
			index := m.chats.Index()
			m.chats.InsertItem(0, recvChat)
			if m.currentChat == "" {
				m.currentChat = recvChat.name
				m.messages.Title = recvChat.name
			} else {
				m.chats.Select(index + 1)
			}
		}
		chatItems[recvChat.name] = append(chatItems[recvChat.name], msg)
		if m.currentChat == recvChat.name {
			m.messages.InsertItem(len(m.messages.Items()), msg)
			m.messages.Select(len(m.messages.Items()) - 1)
		}

	// Key bindings
	case tea.KeyMsg:
		if m.currentView == home {
			if updateHome(&m, &msg, &cmds) {
				return m, tea.Quit
			}
		} else if m.currentView == newDM {
			updateNewDM(&m, &msg, &cmds)
		} else if m.currentView == login {
			if updateLogin(&m, &msg, &cmds) {
				return m, tea.Quit
			}
		} else if m.currentView == help {
			updateHelp(&m, &msg)
		} else if m.currentView == newGC {
			updateNewGroup(&m, &msg, &cmds)
		} else if m.currentView == joinGC {
			updateJoinGroup(&m, &msg, &cmds)
		}

	// Gets terminal Window Size
	case tea.WindowSizeMsg:
		docStyle := lipgloss.NewStyle().Margin(1, 2)
		h, v := docStyle.GetFrameSize()
		m.height = msg.Height - h
		m.width = msg.Width - v
		m.chats.SetSize(m.width/3-1, m.height-4)
		m.messages.SetSize((2*m.width/3)-1, m.height-4)
	}

	return m, tea.Batch(cmds...)
}

func updateHome(m *model, msg *tea.KeyMsg, cmds *[]tea.Cmd) quit {
	inputflag := false
	if m.focus == none {
		switch msg.String() {
		case "q":
			return true
		case "c":
			m.focus = chat
		case "m":
			m.focus = message
		case "n":
			m.currentView = newDM
		case "g":
			m.currentView = newGC
		case "j":
			m.currentView = joinGC
		case "?":
			m.currentView = help
		case "/":
			inputflag = true
			m.focus = input
			m.input.Focus()
		}
	} else {
		switch msg.String() {
		case "esc":
			m.focus = none
			m.input.Blur()
		case "enter":
			if m.focus == chat {
				selectedChat := m.chats.SelectedItem().(Chat).name
				if selectedChat != m.currentChat {
					m.messages.Title = selectedChat
					m.currentChat = selectedChat
					for range len(m.messages.Items()) {
						m.messages.RemoveItem(0)
					}
					for i := range len(chatItems[selectedChat]) {
						m.messages.InsertItem(len(m.messages.Items()), chatItems[selectedChat][i])
					}
				}
			} else if m.focus == input {
				newMsg := Message{Type: "message", From: username, To: m.currentChat, Group: "", IsGroupMsg: false, Content: m.input.Value()}
				if m.chats.SelectedItem().(Chat).isGroupChat {
					newMsg.IsGroupMsg = true
					newMsg.Group = m.currentChat
					newMsg.To = ""
				}
				writeMessages(serverConn, newMsg)
				m.messages.InsertItem(len(m.messages.Items()), newMsg)
				m.messages.Select(len(m.messages.Items()) - 1)
				chatItems[m.currentChat] = append(chatItems[m.currentChat], newMsg)
				m.input.Reset()
			}
		}
	}

	// update focused component
	var cmd tea.Cmd
	if m.focus == chat {
		m.chats, cmd = m.chats.Update(*msg)
		*cmds = append(*cmds, cmd)
	} else if m.focus == message {
		m.messages, cmd = m.messages.Update(*msg)
		*cmds = append(*cmds, cmd)
	} else if m.focus == input && !inputflag {
		m.input, cmd = m.input.Update(*msg)
		*cmds = append(*cmds, cmd)
	}
	return false
}

func updateNewDM(m *model, msg *tea.KeyMsg, cmds *[]tea.Cmd) {
	m.usernameInput.Focus()
	switch msg.String() {
	case "esc":
		m.currentView = home
		m.usernameInput.Reset()
		m.usernameInput.Blur()
	case "enter":
		m.chats.InsertItem(0, Chat{name: m.usernameInput.Value(), isGroupChat: false})
		chatItems[m.usernameInput.Value()] = []Message{}
		m.currentChat = m.usernameInput.Value()
		m.usernameInput.Reset()
		m.usernameInput.Blur()
		m.messages.Title = m.currentChat
		for range len(m.messages.Items()) {
			m.messages.RemoveItem(0)
		}
		m.currentView = home
	}
	var cmd tea.Cmd
	m.usernameInput, cmd = m.usernameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
}

func updateNewGroup(m *model, msg *tea.KeyMsg, cmds *[]tea.Cmd) {
	m.groupnameInput.Focus()
	switch msg.String() {
	case "esc":
		m.currentView = home
		m.groupnameInput.Reset()
		m.groupnameInput.Blur()
	case "enter":
		newMsg := Message{Type: "cmd", From: username, Group: m.groupnameInput.Value(), Content: "create group"}
		writeMessages(serverConn, newMsg)

		m.chats.InsertItem(0, Chat{name: m.groupnameInput.Value(), isGroupChat: true})
		chatItems[m.groupnameInput.Value()] = []Message{}
		m.currentChat = m.groupnameInput.Value()
		m.groupnameInput.Reset()
		m.groupnameInput.Blur()
		m.messages.Title = m.currentChat
		for range len(m.messages.Items()) {
			m.messages.RemoveItem(0)
		}
		m.currentView = home
	}
	var cmd tea.Cmd
	m.groupnameInput, cmd = m.groupnameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
}

func updateJoinGroup(m *model, msg *tea.KeyMsg, cmds *[]tea.Cmd) {
	m.groupnameInput.Focus()
	switch msg.String() {
	case "esc":
		m.currentView = home
		m.groupnameInput.Reset()
		m.groupnameInput.Blur()
	case "enter":
		newMsg := Message{Type: "cmd", From: username, Group: m.groupnameInput.Value(), Content: "join group"}
		writeMessages(serverConn, newMsg)

		m.chats.InsertItem(0, Chat{name: m.groupnameInput.Value(), isGroupChat: true})
		chatItems[m.groupnameInput.Value()] = []Message{}
		m.currentChat = m.groupnameInput.Value()
		m.groupnameInput.Reset()
		m.groupnameInput.Blur()
		m.messages.Title = m.currentChat
		for range len(m.messages.Items()) {
			m.messages.RemoveItem(0)
		}
		m.currentView = home
	}
	var cmd tea.Cmd
	m.groupnameInput, cmd = m.groupnameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
}

func updateLogin(m *model, msg *tea.KeyMsg, cmds *[]tea.Cmd) quit {
	m.usernameInput.Focus()
	switch msg.String() {
	case "q":
		m.usernameInput.Reset()
		m.usernameInput.Blur()
		return true
	case "enter":
		username = m.usernameInput.Value()
		m.usernameInput.Reset()
		m.usernameInput.Blur()
		userLogin()
		m.currentView = home
	}
	var cmd tea.Cmd
	m.usernameInput, cmd = m.usernameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
	return false
}

func updateHelp(m *model, msg *tea.KeyMsg) {
	switch msg.String() {
	case "esc":
		m.currentView = home
	}
}
