package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Message struct {
	Type       string `json:"type"`
	IsGroupMsg bool   `json:"isGroupMsg"`
	Group      string `json:"group"`
	To         string `json:"to"`
	Content    string `json:"content"`
	From       string `json:"From"`
}

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
	name, desc  string
	isGroupChat bool
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
	none component = iota
	input
	chat
	message
)

type view int

const (
	login view = iota
	newDM
	newGC
	joinGC
	home
	help
)

type model struct {
	height         int
	width          int
	chats          list.Model
	messages       list.Model
	input          textinput.Model
	usernameInput  textinput.Model
	groupnameInput textinput.Model
	passwordInput  textinput.Model
	currentChat    string
	focus          component
	currentView    view
}

func (m model) Init() tea.Cmd {
	return nil
}
