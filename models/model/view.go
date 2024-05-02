package model

import (
	"github.com/charmbracelet/lipgloss"
	"signal_client/models"
)

var boxyStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#ebdbb2"))

func (m Model) View() string {
	if m.CurrentView == models.Home {
		return getHomeView(&m)
	} else if m.CurrentView == models.NewDM {
		return getNewDMView(&m)
	} else if m.CurrentView == models.Login {
		return getLoginView(&m)
	} else if m.CurrentView == models.Help {
		return getHelpView(&m)
	} else if m.CurrentView == models.NewGC {
		return getNewGroupView(&m)
	} else if m.CurrentView == models.JoinGC {
		return getJoinGroupView(&m)
	}
	return ""
}

func getHeadView(m *Model) string {
	headHeading := lipgloss.NewStyle().Width(m.Width/2 - 1).Bold(true).Align(lipgloss.Left).Render("Vault V0.1")
	var menuString string
	usernameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b8bb26"))
	if m.CurrentView == models.Home {
		menuString = "[? - help]" + " " + usernameStyle.Render(models.Username)
	} else if m.CurrentView == models.Login {
		menuString = " "
	} else {
		menuString = usernameStyle.Render(models.Username)
	}
	headMenu := lipgloss.NewStyle().Bold(false).Width(m.Width/2 - 1).Align(lipgloss.Right).Render(menuString)
	headView := lipgloss.JoinHorizontal(lipgloss.Center, headHeading, headMenu)
	headStyle := boxyStyle.Width(m.Width).
		PaddingLeft(1).
		Width(m.Width).
		PaddingRight(1).
		Height(1).
		Foreground(lipgloss.Color("#d79921"))
	return headStyle.Render(headView)
}

func getHomeView(m *Model) string {
	inputview := boxyStyle.Width(m.Width).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Render(m.Input.View())
	var chatsview string
	var messagesview string

	if m.CurrentChat != "" {
		chatsview = boxyStyle.Width(m.Width/3 - 1).Bold(false).Height(m.Height - 4).Render(m.Chats.View())
		messagesview = boxyStyle.Width((2 * m.Width / 3) - 1).
			Height(m.Height - 4).
			Bold(false).
			Render(m.Messages.View())
	} else {
		chatWidth := m.Width / 3
		chatsview = boxyStyle.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(chatWidth - 2).Height(m.Height - 4).Render("No Chats")
		messagesview = boxyStyle.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.Width - chatWidth).Height(m.Height - 4).Render("Select a chat")
	}

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m),
		lipgloss.JoinHorizontal(lipgloss.Top, chatsview, messagesview),
		inputview,
	)
}

func getNewDMView(m *Model) string {
	usernameHeadingView := lipgloss.NewStyle().
		Width(m.Width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("Username")
	usernameInputView := boxyStyle.Width(m.Width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.UsernameInput.View())
	newDMView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.Width).
		Height(m.Height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), newDMView)
}

func getNewGroupView(m *Model) string {
	groupNameHeadingView := lipgloss.NewStyle().
		Width(m.Width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("Group Name")
	groupNameInputView := boxyStyle.Width(m.Width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.GroupNameInput.View())
	newGroupView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.Width).
		Height(m.Height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, groupNameHeadingView, groupNameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), newGroupView)
}

func getJoinGroupView(m *Model) string {
	groupNameHeadingView := lipgloss.NewStyle().
		Width(m.Width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("Group Name")
	groupNameInputView := boxyStyle.Width(m.Width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.GroupNameInput.View())
	newGroupView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.Width).
		Height(m.Height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, groupNameHeadingView, groupNameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), newGroupView)
}

func getLoginView(m *Model) string {
	usernameHeadingView := lipgloss.NewStyle().
		Width(m.Width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("Username")
	usernameInputView := boxyStyle.Width(m.Width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.UsernameInput.View())
	// passwordHeadingView := lipgloss.NewStyle().Width(m.Width / 2).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render("Username")
	// passwordInputView := style.Width(m.Width / 2).Align(lipgloss.Left).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render(m.UsernameInput.View())
	loginView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.Width).
		Height(m.Height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))

	// loginView := style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.Width).Height(m.Height - 1).Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView, passwordHeadingView, passwordInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), loginView)
}

func getHelpView(m *Model) string {
	helpHeadingView := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(m.Width / 2).
		Height(1).
		PaddingLeft(1).
		Bold(true).
		Foreground(lipgloss.Color("#11bc7a")).
		Render("Keyboard Shortcuts")
	helplistView := lipgloss.
		NewStyle().
		PaddingLeft(1).
		Align(lipgloss.Left).
		Bold(false).
		Render("\n/ -> Enter Input\nm -> Select Message\nc -> Select Chat\nn -> New Message\nj -> Join Group\ng -> New Group")
	helpView := boxyStyle.Align(lipgloss.Left).
		AlignVertical(lipgloss.Top).
		Width(m.Width).
		Height(m.Height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Top, helpHeadingView, helplistView))
	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), helpView)
}
