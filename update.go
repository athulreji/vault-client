package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	inputflag := false
	switch msg := msg.(type) {

	// Process Incoming Message
	case Message:
		if _, ok := chatItems[msg.From]; !ok {
			index := m.chats.Index()
			m.chats.InsertItem(0, Chat{msg.From, "today"})
			if m.currentChat == "" {
				m.currentChat = msg.From
				m.messages.Title = msg.From
			} else {
				m.chats.Select(index + 1)
			}
		}
		chatItems[msg.From] = append(chatItems[msg.From], msg)
		if m.currentChat == msg.From {
			m.messages.InsertItem(len(m.messages.Items()), msg)
			m.messages.Select(len(m.messages.Items()) - 1)
		}

	// Key bindings
	case tea.KeyMsg:
		if m.currentView == home {
			if m.focus == none {
				switch msg.String() {
				case "q":
					return m, tea.Quit
				case "c":
					m.focus = chat
				case "m":
					m.focus = message
				case "n":
					m.currentView = newDM
					m.usernameInput.Focus()
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
						msg := Message{Type: "message", From: username, To: m.currentChat, Group: "", IsGroupMsg: false, Content: m.input.Value()}
						writeMessages(serverConn, msg)
						m.messages.InsertItem(len(m.messages.Items()), msg)
						m.messages.Select(len(m.messages.Items()) - 1)
						chatItems[m.currentChat] = append(chatItems[m.currentChat], msg)
						m.input.Reset()
					}
				}
			}

			// update focused component
			if m.focus == chat {
				m.chats, cmd = m.chats.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.focus == message {
				m.messages, cmd = m.messages.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.focus == input && !inputflag {
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else if m.currentView == newDM {
			switch msg.String() {
			case "esc":
				m.currentView = home
				m.usernameInput.Reset()
			case "enter":
				m.chats.InsertItem(0, Chat{m.usernameInput.Value(), "today"})
				m.currentChat = m.usernameInput.Value()
				m.usernameInput.Reset()
				m.usernameInput.Blur()
				m.messages.Title = m.currentChat
				for range len(m.messages.Items()) {
					m.messages.RemoveItem(0)
				}
				m.currentView = home
			}
			m.usernameInput, cmd = m.usernameInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.currentView == login {
			m.usernameInput.Focus()
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "enter":
				username = m.usernameInput.Value()
				m.usernameInput.Reset()
				m.usernameInput.Blur()
				userLogin()
				m.currentView = home
			}
			m.usernameInput, cmd = m.usernameInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.currentView == help {
			switch msg.String() {
			case "esc":
				m.currentView = home
			}
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
