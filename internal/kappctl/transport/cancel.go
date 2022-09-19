package transport

import (
	"fmt"
	"github.com/KirillMironov/kappa/pkg/cmdutil"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func (c cmd) cancel() *cobra.Command {
	var apiPath = "/api/v1/deploy"

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

			err = c.client.Do(http.MethodDelete, apiPath, manifest, http.StatusOK)
			if err != nil {
				cmdutil.Exit(err)
			}
		},
	}
}
