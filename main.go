package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	initialInputs = 1
	maxInputs     = 6
	minInputs     = 1
	helpHeight    = 2
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		logger := slog.New(slog.NewJSONHandler(f, nil))
		slog.SetDefault(logger)
		defer f.Close()
	}

	if _, err := tea.NewProgram(newEditor(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
