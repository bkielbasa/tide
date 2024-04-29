package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	initialInputs = 1
	maxInputs     = 6
	minInputs     = 1
	helpHeight    = 2
)

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.ShowLineNumbers = true
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

func main() {
	if _, err := tea.NewProgram(newEditor(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
