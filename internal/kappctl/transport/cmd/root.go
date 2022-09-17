package cmd

import (
	"github.com/spf13/cobra"
)

func (a App) root() *cobra.Command {
	return &cobra.Command{
		Use:   "kappctl",
		Short: "Kappa command-line tool",
	}
}
