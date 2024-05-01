package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle        = lipgloss.NewStyle().PaddingLeft(0).PaddingTop(0).Foreground(lipgloss.Color("#bdae93")).Bold(true)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#689d6a"))
	NoItemStyle       = lipgloss.NewStyle().PaddingLeft(1).PaddingTop(2).Foreground(lipgloss.Color("#689d6a")).Bold(false)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#fabd2f"))
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)
