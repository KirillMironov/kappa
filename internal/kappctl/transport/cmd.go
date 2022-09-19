package transport

import (
	"github.com/KirillMironov/kappa/pkg/cmdutil"
	"github.com/KirillMironov/kappa/pkg/httputil"
	"time"
)

type cmd struct {
	client *httputil.Client
}

func NewCmd(baseApiURL string) *cmd {
	return &cmd{
		client: httputil.NewClient(baseApiURL, time.Second*3),
	}
}

func (c cmd) Execute() {
	root := c.root()
	root.AddCommand(c.deploy(), c.cancel())

	err := root.Execute()
	if err != nil {
		cmdutil.Exit(err)
	}
}
