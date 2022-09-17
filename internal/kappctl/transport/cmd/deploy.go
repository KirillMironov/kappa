package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

var reverseCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a kappa manifest",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manifest, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		client := http.Client{
			Timeout: time.Second * 3,
		}

		resp, err := client.Post("http://localhost:20501/api/v1/deploy", "", manifest)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintln(os.Stderr, resp.Status)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(reverseCmd)
}
