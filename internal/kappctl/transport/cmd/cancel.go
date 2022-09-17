package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func (a App) cancel() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a kappa deployment",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			manifest, err := os.Open(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer manifest.Close()

			resp, err := a.requester.Do(http.MethodDelete, "/api/v1/deploy", manifest)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println(resp)
		},
	}
}
