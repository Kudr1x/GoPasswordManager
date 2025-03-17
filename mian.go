package main

import (
	"goPasswordManager/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
