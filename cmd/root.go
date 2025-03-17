package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "pm",
	Short: "Deterministic Password Manager",
	Long:  `A deterministic password manager using master password and service name`,
}
