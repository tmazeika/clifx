package main

import (
	_ "github.com/inconshreveable/mousetrap"
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
