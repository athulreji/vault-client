package main

import "github.com/charmbracelet/lipgloss"

var style = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#57288d"))

func (m model) View() string {
	if m.currentView == home {
		return getHomeView(&m)
	} else if m.currentView == newDM {
		return getNewDMView(&m)
	} else if m.currentView == login {
		return getLoginView(&m)
	}
	return ""
}

func getHeadView(m *model) string {
	return style.Width(m.width).
		Align(lipgloss.Left).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#11bc7a")).
		Render(" Vault V0.1")
}

func getHomeView(m *model) string {
	inputview := style.Width(m.width).
		Align(lipgloss.Left).
		Height(1).
		Foreground(lipgloss.Color("#11bc7a")).
		Render(m.input.View())
	chatsview := style.Width(m.width/3 - 1).Height(m.height - 4).Render(m.chats.View())
	var messagesview string

	if m.currentChat != "" {
		messagesview = style.Width((2 * m.width / 3) - 1).
			Height(m.height - 4).
			Render(m.messages.View())
	} else {
		messagesview = style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width((2 * m.width / 3) - 1).Height(m.height - 4).Render("Select a chat")
	}

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m),
		lipgloss.JoinHorizontal(lipgloss.Top, chatsview, messagesview),
		inputview,
	)
}

func getNewDMView(m *model) string {
	usernameHeadingView := lipgloss.NewStyle().
		Width(m.width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#11bc7a")).
		Render("Username")
	usernameInputView := style.Width(m.width / 2).
		Align(lipgloss.Left).
		Height(1).
		Foreground(lipgloss.Color("#11bc7a")).
		Render(m.usernameInput.View())
	newDMview := style.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), newDMview)
}

func getLoginView(m *model) string {
	usernameHeadingView := lipgloss.NewStyle().
		Width(m.width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#11bc7a")).
		Render("Username")
	usernameInputView := style.Width(m.width / 2).
		Align(lipgloss.Left).
		Height(1).
		Foreground(lipgloss.Color("#11bc7a")).
		Render(m.usernameInput.View())
	// passwordHeadingView := lipgloss.NewStyle().Width(m.width / 2).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render("Username")
	// passwordInputView := style.Width(m.width / 2).Align(lipgloss.Left).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render(m.usernameInput.View())
	loginView := style.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))

		// loginView := style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.width).Height(m.height - 1).Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView, passwordHeadingView, passwordInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), loginView)
}
