package cmd

import (
	"fmt"
	"io"
	"os"
)

type App struct {
	requester requester
}

type requester interface {
	Do(method string, path string, body io.Reader) (respBody string, err error)
}

func NewApp(requester requester) *App {
	return &App{
		requester: requester,
	}
}

func (a App) Execute() {
	root := a.root()
	root.AddCommand(a.deploy(), a.cancel())

	err := root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
