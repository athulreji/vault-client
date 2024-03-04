package main

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	var style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#57288d"))

	headview := style.Width(m.width).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render(" Vault V0.1")
	inputview := style.Width(m.width).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render(m.input.View())
	chatsview := style.Width(m.width/3 - 1).Height(m.height - 4).Render(m.chats.View())
	var messagesview string

	if m.currentChat != "" {
		messagesview = style.Width((2 * m.width / 3) - 1).Height(m.height - 4).Render(m.messages.View())
	} else {
		messagesview = style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width((2 * m.width / 3) - 1).Height(m.height - 4).Render("Select a chat")
	}

	if m.currentView == home {
		return lipgloss.JoinVertical(lipgloss.Top, headview,
			lipgloss.JoinHorizontal(lipgloss.Top, chatsview, messagesview),
			inputview,
		)
	} else if m.currentView == newDM {
		usernameHeadingView := lipgloss.NewStyle().Width(m.width / 2).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render("Username")
		usernameInputView := style.Width(m.width / 2).Align(lipgloss.Left).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render(m.usernameInput.View())
		newDMview := style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.width).Height(m.height - 1).Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))

		return lipgloss.JoinVertical(lipgloss.Top, headview, newDMview)
	} else if m.currentView == login {
		usernameHeadingView := lipgloss.NewStyle().Width(m.width / 2).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render("Username")
		usernameInputView := style.Width(m.width / 2).Align(lipgloss.Left).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render(m.usernameInput.View())
		//passwordHeadingView := lipgloss.NewStyle().Width(m.width / 2).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render("Username")
		//passwordInputView := style.Width(m.width / 2).Align(lipgloss.Left).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render(m.usernameInput.View())
		loginView := style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.width).Height(m.height - 1).Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))
		//loginView := style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.width).Height(m.height - 1).Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView, passwordHeadingView, passwordInputView))

		return lipgloss.JoinVertical(lipgloss.Top, headview, loginView)
	}
	return ""
}
