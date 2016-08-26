package main

import (
	"fmt"
	"github.com/bionicrm/clifx/cmd"
	_ "github.com/inconshreveable/mousetrap"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
