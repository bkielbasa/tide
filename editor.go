package main

import (
	"context"
	"fmt"
	"github.com/bkielbasa/tide/cmdinput"
	"github.com/bkielbasa/tide/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type editor struct {
	width   int
	height  int
	buffers []buffer
	focus   int
	mode    editorMode

	normalModeCommands []command
	isEnteringCommand  bool
	commandInput       cmdinput.Model
}

type editorMode string

const (
	modeNormal editorMode = "NORMAL"
	modeInsert editorMode = "INSERT"
)

func newEditor() *editor {
	m := editor{
		buffers: make([]buffer, initialInputs),
		mode:    modeNormal,
		normalModeCommands: []command{
			cmdQuit{},
		},
		commandInput: cmdinput.New(),
	}
	for i := 0; i < initialInputs; i++ {
		m.buffers[i] = newBuffer()
	}
	m.buffers[m.focus].Focus()
	return &m
}

func (m editor) Init() tea.Cmd {
	return textarea.Blink
}

func (m editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var updateBuffers bool
	var updateInput bool

	switch m.mode {
	case modeNormal:
		var cmd tea.Cmd
		updateBuffers, updateInput, m, cmd = m.updateNormalMode(msg)
		cmds = append(cmds, cmd)

	case modeInsert:
		var cmd tea.Cmd
		updateBuffers, m, cmd = m.updateInsertMode(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.sizeBuffers()

	if updateInput {
		newModel, cmd := m.commandInput.Update(msg)
		m.commandInput = newModel
		cmds = append(cmds, cmd)
	}

	if updateBuffers {
		for i := range m.buffers {
			newModel, cmd := m.buffers[i].Update(msg)
			m.buffers[i] = newModel
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *editor) sizeBuffers() {
	for i := range m.buffers {
		m.buffers[i].SetWidth(m.width / len(m.buffers))
		m.buffers[i].SetHeight(m.height - helpHeight)
	}
}

func (m editor) updateNormalMode(msg tea.Msg) (bool, bool, editor, tea.Cmd) {
	var cmds []tea.Cmd

	updateInputs := true
	updateCommandInput := false

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.isEnteringCommand {
			switch msg.String() {
			case "esc":
				m.buffers[0].Focus()
				m.commandInput.Blur()
				m.isEnteringCommand = false
				m.commandInput.Reset()
				updateCommandInput = false
			case "enter":
				currValue := m.commandInput.Value()
				if currValue == "" {
					updateCommandInput = false
					break
				}

				val := strings.Split(currValue, " ")

				// TODO: add proper command parsing
				for _, cmd := range m.normalModeCommands {
					if val[0] == cmd.Short() || val[0] == cmd.FullName() {
						c, err := cmd.Exec(context.Background(), val[1:]...)
						if err != nil {
							fmt.Print(err)
						}

						return false, true, m, c
					}
				}
			default:
				updateCommandInput = true
			}
		} else {
			switch msg.String() {
			case "k":
				m.buffers[0].CursorUp()
				updateInputs = false
			case "h":
				m.buffers[0].CharacterLeft()
				updateInputs = false
			case "l":
				m.buffers[0].CharacterRight()
				updateInputs = false
			case "j":
				m.buffers[0].CursorDown()
				updateInputs = false
			case ":":
				m.buffers[0].Blur()
				m.commandInput.SetValue("")
				cmds = append(cmds, m.commandInput.Focus())
				m.isEnteringCommand = true
			case "i":
				m.mode = modeInsert
				updateInputs = false
			default:
				updateInputs = false
			}
		}

	}

	return updateInputs, updateCommandInput, m, tea.Batch(cmds...)
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
	}

	return updateInputs, m, tea.Batch(cmds...)
}

func (m editor) View() string {
	var views []string
	for i := range m.buffers {
		views = append(views, m.buffers[i].View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, views...), m.commandInput.View()) + "\n" + string(m.mode)
}
