package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"signal_client/models"
	"signal_client/models/chat"
	"signal_client/models/message"
	"signal_client/services"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	// Process Incoming Message
	case message.Message:
		recvChat := chat.Chat{}
		if msg.IsGroupMsg {
			recvChat.IsGroupChat = true
			recvChat.Name = msg.Group
		} else {
			recvChat.IsGroupChat = false
			recvChat.Name = msg.From
		}
		if _, ok := models.ChatItems[recvChat.Name]; !ok {
			index := m.Chats.Index()
			m.Chats.InsertItem(0, recvChat)
			if m.CurrentChat == "" {
				m.CurrentChat = recvChat.Name
				m.Messages.Title = recvChat.Name
			} else {
				m.Chats.Select(index + 1)
			}
		}
		models.ChatItems[recvChat.Name] = append(models.ChatItems[recvChat.Name], msg)
		if m.CurrentChat == recvChat.Name {
			m.Messages.InsertItem(len(m.Messages.Items()), msg)
			m.Messages.Select(len(m.Messages.Items()) - 1)
		}

	// Key bindings
	case tea.KeyMsg:
		if m.CurrentView == models.Home {
			if UpdateHome(&m, &msg, &cmds) {
				return m, tea.Quit
			}
		} else if m.CurrentView == models.NewDM {
			UpdateNewDM(&m, &msg, &cmds)
		} else if m.CurrentView == models.Login {
			if UpdateLogin(&m, &msg, &cmds) {
				return m, tea.Quit
			}
		} else if m.CurrentView == models.Help {
			UpdateHelp(&m, &msg)
		} else if m.CurrentView == models.NewGC {
			UpdateNewGroup(&m, &msg, &cmds)
		} else if m.CurrentView == models.JoinGC {
			UpdateJoinGroup(&m, &msg, &cmds)
		}

	// Gets terminal Window Size
	case tea.WindowSizeMsg:
		docStyle := lipgloss.NewStyle().Margin(1, 2)
		h, v := docStyle.GetFrameSize()
		m.Height = msg.Height - h
		m.Width = msg.Width - v
		m.Chats.SetSize(m.Width/3-1, m.Height-4)
		m.Messages.SetSize((2*m.Width/3)-1, m.Height-4)
	}

	return m, tea.Batch(cmds...)
}
func UpdateHome(m *Model, msg *tea.KeyMsg, cmds *[]tea.Cmd) bool {
	inputflag := false
	if m.Focus == models.None {
		switch msg.String() {
		case "q":
			return true
		case "c":
			m.Focus = models.Chat
		case "m":
			m.Focus = models.Message
		case "n":
			m.CurrentView = models.NewDM
		case "g":
			m.CurrentView = models.NewGC
		case "j":
			m.CurrentView = models.JoinGC
		case "?":
			m.CurrentView = models.Help
		case "/":
			inputflag = true
			m.Focus = models.Input
			m.Input.Focus()
		}
	} else {
		switch msg.String() {
		case "esc":
			m.Focus = models.None
			m.Input.Blur()
		case "enter":
			if m.Focus == models.Chat {
				selectedChat := m.Chats.SelectedItem().(chat.Chat).Name
				if selectedChat != m.CurrentChat {
					m.Messages.Title = selectedChat
					m.CurrentChat = selectedChat
					for range len(m.Messages.Items()) {
						m.Messages.RemoveItem(0)
					}
					for i := range len(models.ChatItems[selectedChat]) {
						m.Messages.InsertItem(len(m.Messages.Items()), models.ChatItems[selectedChat][i])
					}
				}
			} else if m.Focus == models.Input {
				newMsg := message.Message{Type: "message", From: models.Username, To: m.CurrentChat, Group: "", IsGroupMsg: false, Content: m.Input.Value()}
				if m.Chats.SelectedItem().(chat.Chat).IsGroupChat {
					newMsg.IsGroupMsg = true
					newMsg.Group = m.CurrentChat
					newMsg.To = ""
				}
				services.WriteMessages(models.ServerConn, newMsg)
				m.Messages.InsertItem(len(m.Messages.Items()), newMsg)
				m.Messages.Select(len(m.Messages.Items()) - 1)
				models.ChatItems[m.CurrentChat] = append(models.ChatItems[m.CurrentChat], newMsg)
				m.Input.Reset()
			}
		}
	}

	// update focused component
	var cmd tea.Cmd
	if m.Focus == models.Chat {
		m.Chats, cmd = m.Chats.Update(*msg)
		*cmds = append(*cmds, cmd)
	} else if m.Focus == models.Message {
		m.Messages, cmd = m.Messages.Update(*msg)
		*cmds = append(*cmds, cmd)
	} else if m.Focus == models.Input && !inputflag {
		m.Input, cmd = m.Input.Update(*msg)
		*cmds = append(*cmds, cmd)
	}
	return false
}

func UpdateNewDM(m *Model, msg *tea.KeyMsg, cmds *[]tea.Cmd) {
	m.UsernameInput.Focus()
	switch msg.String() {
	case "esc":
		m.CurrentView = models.Home
		m.UsernameInput.Reset()
		m.UsernameInput.Blur()
	case "enter":
		m.Chats.InsertItem(0, chat.Chat{Name: m.UsernameInput.Value(), IsGroupChat: false})
		models.ChatItems[m.UsernameInput.Value()] = []message.Message{}
		m.CurrentChat = m.UsernameInput.Value()
		m.UsernameInput.Reset()
		m.UsernameInput.Blur()
		m.Messages.Title = m.CurrentChat
		for range len(m.Messages.Items()) {
			m.Messages.RemoveItem(0)
		}
		m.CurrentView = models.Home
	}
	var cmd tea.Cmd
	m.UsernameInput, cmd = m.UsernameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
}

func UpdateNewGroup(m *Model, msg *tea.KeyMsg, cmds *[]tea.Cmd) {
	m.GroupNameInput.Focus()
	switch msg.String() {
	case "esc":
		m.CurrentView = models.Home
		m.GroupNameInput.Reset()
		m.GroupNameInput.Blur()
	case "enter":
		newMsg := message.Message{Type: "cmd", From: models.Username, Group: m.GroupNameInput.Value(), Content: "create group"}
		services.WriteMessages(models.ServerConn, newMsg)

		m.Chats.InsertItem(0, chat.Chat{Name: m.GroupNameInput.Value(), IsGroupChat: true})
		models.ChatItems[m.GroupNameInput.Value()] = []message.Message{}
		m.CurrentChat = m.GroupNameInput.Value()
		m.GroupNameInput.Reset()
		m.GroupNameInput.Blur()
		m.Messages.Title = m.CurrentChat
		for range len(m.Messages.Items()) {
			m.Messages.RemoveItem(0)
		}
		m.CurrentView = models.Home
	}
	var cmd tea.Cmd
	m.GroupNameInput, cmd = m.GroupNameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
}

func UpdateJoinGroup(m *Model, msg *tea.KeyMsg, cmds *[]tea.Cmd) {
	m.GroupNameInput.Focus()
	switch msg.String() {
	case "esc":
		m.CurrentView = models.Home
		m.GroupNameInput.Reset()
		m.GroupNameInput.Blur()
	case "enter":
		newMsg := message.Message{Type: "cmd", From: models.Username, Group: m.GroupNameInput.Value(), Content: "join group"}
		services.WriteMessages(models.ServerConn, newMsg)

		m.Chats.InsertItem(0, chat.Chat{Name: m.GroupNameInput.Value(), IsGroupChat: true})
		models.ChatItems[m.GroupNameInput.Value()] = []message.Message{}
		m.CurrentChat = m.GroupNameInput.Value()
		m.GroupNameInput.Reset()
		m.GroupNameInput.Blur()
		m.Messages.Title = m.CurrentChat
		for range len(m.Messages.Items()) {
			m.Messages.RemoveItem(0)
		}
		m.CurrentView = models.Home
	}
	var cmd tea.Cmd
	m.GroupNameInput, cmd = m.GroupNameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
}

func UpdateLogin(m *Model, msg *tea.KeyMsg, cmds *[]tea.Cmd) bool {
	m.UsernameInput.Focus()
	switch msg.String() {
	case "q":
		m.UsernameInput.Reset()
		m.UsernameInput.Blur()
		return true
	case "enter":
		models.Username = m.UsernameInput.Value()
		m.UsernameInput.Reset()
		m.UsernameInput.Blur()
		services.UserLogin()
		m.CurrentView = models.Home
	}
	var cmd tea.Cmd
	m.UsernameInput, cmd = m.UsernameInput.Update(*msg)
	*cmds = append(*cmds, cmd)
	return false
}

func UpdateHelp(m *Model, msg *tea.KeyMsg) {
	switch msg.String() {
	case "esc":
		m.CurrentView = models.Home
	}
}
