package deligates

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"signal_client/models/chat"
	message2 "signal_client/models/message"
	"signal_client/style"
	"strings"
)

type ChatItemDelegate struct{}
type MessageItemDelegate struct{}

func (d ChatItemDelegate) Height() int                             { return 1 }
func (d ChatItemDelegate) Spacing() int                            { return 0 }
func (d ChatItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ChatItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(chat.Chat)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Name)

	fn := style.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

func (d MessageItemDelegate) Height() int                             { return 1 }
func (d MessageItemDelegate) Spacing() int                            { return 0 }
func (d MessageItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d MessageItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(message2.Message)
	if !ok {
		return
	}

	head := lipgloss.NewStyle().Bold(true).Render(i.From)

	str := fmt.Sprintf("%s: %s", head, i.Content)

	fn := style.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}
