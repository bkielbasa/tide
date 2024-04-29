package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type buffer struct {
	name string
	text textarea.Model
}

func newBuffer() buffer {
	t := textarea.New()
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.ShowLineNumbers = true
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))

	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	t.Cursor.Style = cursorStyle

	b := buffer{
		text: t,
		name: "UNNAMED",
	}

	b.Open("main.go")

	return b
}

func (b *buffer) Open(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("buffer cannot open file: %w", err)
	}

	cont, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("buffer cannot read file: %w", err)
	}

	b.text.SetValue(string(cont))

	return nil
}

func (b *buffer) CursorUp() {
	b.text.CursorUp()
}

func (b *buffer) CursorDown() {
	b.text.CursorDown()
}

func (b *buffer) SetSize(width, height int) {
	b.text.SetWidth(width)
	b.text.SetHeight(height)
}

func (b *buffer) Focus() tea.Cmd {
	return b.text.Focus()
}

func (b *buffer) Blur() {
	b.text.Blur()
}

func (b buffer) Update(msg tea.Msg) (buffer, tea.Cmd) {
	m, cmd := b.text.Update(msg)
	b.text = m

	return b, cmd
}

func (b buffer) View() string {
	return b.text.View()
}

func (b buffer) SetWidth(width int) buffer {
	b.text.SetWidth(width)
	return b
}

func (b buffer) SetHeight(height int) buffer {
	b.text.SetHeight(height)
	return b
}
