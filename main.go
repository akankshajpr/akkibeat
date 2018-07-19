package main

import (
	"os"

	"github.com/akankshajpr/akkibeat/cmd"

	_ "github.com/akankshajpr/akkibeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
