package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

func NewDefaultCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kappctl",
		Short: "Kappa command-line tool",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Hello, world!")
		},
	}
}
