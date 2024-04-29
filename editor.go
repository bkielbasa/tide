package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type editor struct {
	width  int
	height int
	inputs []buffer
	focus  int
	mode   editorMode
}

type editorMode string

const (
	modeNormal editorMode = "NORMAL"
	modeInsert editorMode = "INSERT"
)

func newEditor() *editor {
	m := editor{
		inputs: make([]buffer, initialInputs),
		mode:   modeNormal,
	}
	for i := 0; i < initialInputs; i++ {
		m.inputs[i] = newBuffer()
	}
	m.inputs[m.focus].Focus()
	return &m
}

func (m editor) Init() tea.Cmd {
	return textarea.Blink
}

func (m editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	updateInputs := true

	switch m.mode {
	case modeNormal:
		var cmd tea.Cmd
		updateInputs, m, cmd = m.updateNormalMode(msg)
		cmds = append(cmds, cmd)

	case modeInsert:
		var cmd tea.Cmd
		updateInputs, m, cmd = m.updateInsertMode(msg)
		cmds = append(cmds, cmd)
	}

	m.sizeInputs()

	if updateInputs {
		// Update all textareas
		for i := range m.inputs {
			newModel, cmd := m.inputs[i].Update(msg)
			m.inputs[i] = newModel
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *editor) sizeInputs() {
	for i := range m.inputs {
		m.inputs[i].SetWidth(m.width / len(m.inputs))
		m.inputs[i].SetHeight(m.height - helpHeight)
	}
}

func (m editor) updateNormalMode(msg tea.Msg) (bool, editor, tea.Cmd) {
	var cmds []tea.Cmd

	updateInputs := true

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			m.inputs[0].CursorUp()
		case "j":
			m.inputs[0].CursorDown()
		case "q":
			return updateInputs, m, tea.Quit
		case "i":
			m.mode = modeInsert
			updateInputs = false
		default:
			updateInputs = false
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	return updateInputs, m, tea.Batch(cmds...)
}

func (m editor) updateInsertMode(msg tea.Msg) (bool, editor, tea.Cmd) {
	var cmds []tea.Cmd
	updateInputs := true

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = modeNormal
			updateInputs = false
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	return updateInputs, m, tea.Batch(cmds...)
}

func (m editor) View() string {
	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + string(m.mode)
}
