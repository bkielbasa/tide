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

type bufferOpener interface {
	OpenBuffer(fileName string) error
}

type cmdOpenFiles struct {
	opener bufferOpener
}

func (c cmdOpenFiles) FullName() string {
	return "open"
}

func (c cmdOpenFiles) Short() string {
	return "o"
}

func (c cmdOpenFiles) Exec(ctx context.Context, args ...string) (tea.Cmd, error) {
	for _, f := range args {
		if err := c.opener.OpenBuffer(f); err != nil {
			return commandCloseCommandInput, err
		}
	}

	return commandCloseCommandInput, nil
}
