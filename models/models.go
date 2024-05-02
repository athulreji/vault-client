package models

import message2 "signal_client/models/message"

var ChatItems = make(map[string][]message2.Message)

type Component int
type View int

const (
	None Component = iota
	Input
	Chat
	Message
)

const (
	Login View = iota
	NewDM
	NewGC
	JoinGC
	Home
	Help
)
