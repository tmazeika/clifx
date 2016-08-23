package main

import (
	"github.com/bionicrm/clifx/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
