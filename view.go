package main

import (
	"github.com/charmbracelet/lipgloss"
)

var boxyStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#ebdbb2"))

func (m model) View() string {
	if m.currentView == home {
		return getHomeView(&m)
	} else if m.currentView == newDM {
		return getNewDMView(&m)
	} else if m.currentView == login {
		return getLoginView(&m)
	} else if m.currentView == help {
		return getHelpView(&m)
	} else if m.currentView == newGC {
		return getNewGroupView(&m)
	} else if m.currentView == joinGC {
		return getJoinGroupView(&m)
	} else if m.currentView == sendFile {
		return getFilePickerView(&m)
	}
	return ""
}

func getHeadView(m *model) string {
	headHeading := lipgloss.NewStyle().Width(m.width/2 - 1).Bold(true).Align(lipgloss.Left).Render("Vault V0.1")
	var menuString string
	usernameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b8bb26"))
	if m.currentView == home {
		menuString = "[? - help]" + " " + usernameStyle.Render(username)
	} else if m.currentView == login {
		menuString = " "
	} else {
		menuString = usernameStyle.Render(username)
	}
	headMenu := lipgloss.NewStyle().Bold(false).Width(m.width/2 - 1).Align(lipgloss.Right).Render(menuString)
	headView := lipgloss.JoinHorizontal(lipgloss.Center, headHeading, headMenu)
	headStyle := boxyStyle.Width(m.width).
		PaddingLeft(1).
		Width(m.width).
		PaddingRight(1).
		Height(1).
		Foreground(lipgloss.Color("#d79921"))
	return headStyle.Render(headView)
}

func getHomeView(m *model) string {
	inputview := boxyStyle.Width(m.width).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Render(m.input.View())
	var chatsview string
	var messagesview string

	if m.currentChat != "" {
		chatsview = boxyStyle.Width(m.width/3 - 1).Bold(false).Height(m.height - 4).Render(m.chats.View())
		messagesview = boxyStyle.Width((2 * m.width / 3) - 1).
			Height(m.height - 4).
			Bold(false).
			Render(m.messages.View())
	} else {
		chatWidth := m.width / 3
		chatsview = boxyStyle.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(chatWidth - 2).Height(m.height - 4).Render("No chats")
		messagesview = boxyStyle.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.width - chatWidth).Height(m.height - 4).Render("Select a chat")
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
		Foreground(lipgloss.Color("#fe8019")).
		Render("Username")
	usernameInputView := boxyStyle.Width(m.width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.usernameInput.View())
	newDMView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), newDMView)
}

func getNewGroupView(m *model) string {
	groupNameHeadingView := lipgloss.NewStyle().
		Width(m.width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("Group Name")
	groupNameInputView := boxyStyle.Width(m.width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.groupnameInput.View())
	newGroupView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, groupNameHeadingView, groupNameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), newGroupView)
}

func getJoinGroupView(m *model) string {
	groupNameHeadingView := lipgloss.NewStyle().
		Width(m.width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("Group Name")
	groupNameInputView := boxyStyle.Width(m.width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.groupnameInput.View())
	newGroupView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, groupNameHeadingView, groupNameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), newGroupView)
}

func getLoginView(m *model) string {
	usernameHeadingView := lipgloss.NewStyle().
		Width(m.width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("Username")
	usernameInputView := boxyStyle.Width(m.width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.usernameInput.View())
	// passwordHeadingView := lipgloss.NewStyle().Width(m.width / 2).Height(1).Bold(true).Foreground(lipgloss.Color("#11bc7a")).Render("Username")
	// passwordInputView := style.Width(m.width / 2).Align(lipgloss.Left).Height(1).Foreground(lipgloss.Color("#11bc7a")).Render(m.usernameInput.View())
	loginView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView))

		// loginView := style.Align(lipgloss.Center).AlignVertical(lipgloss.Center).Width(m.width).Height(m.height - 1).Render(lipgloss.JoinVertical(lipgloss.Center, usernameHeadingView, usernameInputView, passwordHeadingView, passwordInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), loginView)
}

func getHelpView(m *model) string {
	helpHeadingView := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(m.width / 2).
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
		Render("\n/ -> Enter Input\nm -> Select Message\nc -> Select Chat\nn -> New Message\nj -> Join Group\ng -> New Group\nf -> Send File\no -> Open file")
	helpView := boxyStyle.Align(lipgloss.Left).
		AlignVertical(lipgloss.Top).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Top, helpHeadingView, helplistView))
	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), helpView)
}

func getFilePickerView(m *model) string {
	filenameHeadingView := lipgloss.NewStyle().
		Width(m.width / 2).
		Height(1).
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Render("File Name")
	filenameInputView := boxyStyle.Width(m.width / 2).
		Align(lipgloss.Left).
		Height(1).
		Bold(false).
		Foreground(lipgloss.Color("#83a598")).
		Render(m.filenameInput.View())
	filepickerView := boxyStyle.Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(m.width).
		Height(m.height - 1).
		Render(lipgloss.JoinVertical(lipgloss.Center, filenameHeadingView, filenameInputView))

	return lipgloss.JoinVertical(lipgloss.Top, getHeadView(m), filepickerView)

	// if m.quitting {
	// 	return ""
	// }
	// var s strings.Builder
	// s.WriteString("\n  ")
	// if m.err != nil {
	// 	s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	// } else if m.selectedFile == "" {
	// 	s.WriteString("Pick a file:")
	// } else {
	// 	s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	// }
	// s.WriteString("\n\n" + m.filepicker.View() + "\n")
	// return s.String()
}
