package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

type command interface {
	FullName() string
	Short() string
	Exec(ctx context.Context, args ...string) (tea.Cmd, error)
}

type cmdQuit struct {
}

func (cmdQuit) FullName() string {
	return "quit"
}

func (cmdQuit) Short() string {
	return "q"
}

func (cmdQuit) Exec(ctx context.Context, args ...string) (tea.Cmd, error) {
	return tea.Quit, nil
}
