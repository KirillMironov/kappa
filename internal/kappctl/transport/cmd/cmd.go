package cmd

import (
	"fmt"
	"io"
	"os"
)

type Cmd struct {
	requester requester
}

type requester interface {
	Do(method string, path string, body io.Reader) (respBody string, err error)
}

func NewCmd(requester requester) *Cmd {
	return &Cmd{
		requester: requester,
	}
}

func (c Cmd) Execute() {
	root := c.root()
	root.AddCommand(c.deploy(), c.cancel())

	err := root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
