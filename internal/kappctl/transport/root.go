package transport

import (
	"github.com/spf13/cobra"
)

func (c cmd) root() *cobra.Command {
	return &cobra.Command{
		Use:   "kappctl",
		Short: "Kappa command-line tool",
	}
}
