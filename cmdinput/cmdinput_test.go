package cmdinput_test

import (
	"testing"

	"github.com/bkielbasa/tide/cmdinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCmdInput_Update(t *testing.T) {
	m := cmdinput.New()

	t.Run("adding new character", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
		m.Focus()
		m.SetValue("a")
		m, _ = m.Update(msg)
		if m.Value() != "ab" {
			t.Errorf("expected %q but received %q", "ab", m.Value())
		}
	})

	t.Run("remove a character", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyBackspace}
		m.Focus()
		m.SetValue("a")
		m, _ = m.Update(msg)
		if m.Value() != "" {
			t.Errorf("expected %q but received %q", "", m.Value())
		}
	})

	t.Run("move one character left and right", func(t *testing.T) {
		m.Focus()
		m.SetValue("abcd")

		// the cursor should point to the last character
		col := m.CursorCol()
		if col != 4 {
			t.Errorf("the cursor should point to pos 4 but points to %d", col)
		}

		msg := tea.KeyMsg{Type: tea.KeyLeft}
		m, _ = m.Update(msg)
		if m.CursorCol() != col-1 {
			t.Errorf("the cursor should point to pos %d but points to %d", col-1, m.CursorCol())
		}

		msg = tea.KeyMsg{Type: tea.KeyRight}
		m, _ = m.Update(msg)
		if m.CursorCol() != col {
			t.Errorf("the cursor should point to pos %d but points to %d", col, m.CursorCol())
		}
	})
}
