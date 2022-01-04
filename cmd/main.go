package main

import (
	"log"
	"os"
)

var (
	version string = "dev"
)

func main() {
	c := newCoverCommand(version)
	if err := c.Run(os.Args[1:]); err != nil {
		log.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
