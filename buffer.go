package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bkielbasa/tide/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type buffer struct {
	name string
	text textarea.Model

	parser *sitter.Parser
}

func newBuffer() buffer {
	t := textarea.New()
	t.CharLimit = 0
	t.ShowLineNumbers = true

	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	t.Cursor.Style = cursorStyle

	b := buffer{
		text:   t,
		name:   "UNNAMED",
		parser: sitter.NewParser(),
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
	b.text.MoveToBegin()

	b.parser.SetLanguage(golang.GetLanguage())

	return nil
}

func (b *buffer) CursorUp() {
	b.text.CursorUp()
}

func (b *buffer) CharacterLeft() {
	b.text.CharacterLeft(false)
}

func (b *buffer) CharacterRight() {
	b.text.CharacterRight()
}

func (b *buffer) CursorDown() {
	b.text.CursorDown()
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

func (b *buffer) SetWidth(width int) {
	b.text.SetWidth(width)
}

func (b *buffer) SetHeight(height int) {
	b.text.SetHeight(height)
}
