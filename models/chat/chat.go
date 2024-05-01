package chat

type Chat struct {
	Name, Desc  string
	IsGroupChat bool
}

func (item Chat) Title() string {
	return item.Name
}

func (item Chat) Description() string {
	return item.Desc
}

func (item Chat) FilterValue() string {
	return item.Name
}
