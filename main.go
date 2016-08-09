package main

import (
	"github.com/bionicrm/clifx/cmd"
	"log"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
