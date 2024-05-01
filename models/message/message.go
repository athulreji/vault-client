package message

type Message struct {
	Type       string `json:"type"`
	IsGroupMsg bool   `json:"isGroupMsg"`
	Group      string `json:"group"`
	To         string `json:"to"`
	Content    string `json:"content"`
	From       string `json:"from"`
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
