package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chatItemDelegate struct{}

func (d chatItemDelegate) Height() int                             { return 1 }
func (d chatItemDelegate) Spacing() int                            { return 0 }
func (d chatItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d chatItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Chat)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type messageItemDelegate struct{}

func (d messageItemDelegate) Height() int                             { return 1 }
func (d messageItemDelegate) Spacing() int                            { return 0 }
func (d messageItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d messageItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Message)
	if !ok {
		return
	}

	head := lipgloss.NewStyle().Bold(true).Render(i.From)

	str := fmt.Sprintf("%s: %s", head, i.Content)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
