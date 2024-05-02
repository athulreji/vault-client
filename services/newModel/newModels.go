package newModel

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"signal_client/models/deligates"
	"signal_client/models/model"
	"signal_client/style"
)

// NewModel initializes the UI and returns a new instance of model.
func NewModel() model.Model {
	chatList := list.New([]list.Item{}, deligates.ChatItemDelegate{}, 0, 0)
	chatList.SetShowHelp(false)
	chatList.Styles = list.Styles{
		TitleBar:        style.TitleStyle,
		NoItems:         style.NoItemStyle,
		PaginationStyle: style.PaginationStyle,
	}
	chatList.SetShowFilter(false)
	chatList.SetShowStatusBar(false)

	input := textinput.New()
	usernameInput := textinput.New()
	passwordInput := textinput.New()
	groupNameInput := textinput.New()

	messageList := list.New([]list.Item{}, deligates.MessageItemDelegate{}, 0, 0)
	messageList.Styles = list.Styles{
		NoItems: lipgloss.NewStyle().PaddingLeft(2).PaddingTop(2),
	}
	messageList.SetShowHelp(false)
	messageList.SetShowTitle(false)
	messageList.SetShowStatusBar(false)
	messageList.SetShowFilter(false)

	return model.Model{
		Chats:          chatList,
		Messages:       messageList,
		Input:          input,
		UsernameInput:  usernameInput,
		PasswordInput:  passwordInput,
		GroupNameInput: groupNameInput,
	}
}
