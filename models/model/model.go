package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"signal_client/models"
)

type Model struct {
	Height         int
	Width          int
	Chats          list.Model
	Messages       list.Model
	Input          textinput.Model
	UsernameInput  textinput.Model
	GroupNameInput textinput.Model
	PasswordInput  textinput.Model
	CurrentChat    string
	Focus          models.Component
	CurrentView    models.View
}

func (m Model) Init() tea.Cmd {
	return nil
}
