package cmdutil

import (
	"fmt"
	"os"
)

func Exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
