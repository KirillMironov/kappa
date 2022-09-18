package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func (c Cmd) deploy() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a kappa manifest",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			manifest, err := os.Open(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer manifest.Close()

			resp, err := c.requester.Do(http.MethodPost, "/api/v1/deploy", manifest)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println(resp)
		},
	}
}
