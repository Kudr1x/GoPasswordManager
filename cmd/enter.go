package cmd

import (
	"crypto/subtle"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"goPasswordManager/internal/crypto"
	"goPasswordManager/internal/storage"
)

var enterCmd = &cobra.Command{
	Use:   "enter [name]",
	Short: "Access password storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		storageName := args[0]

		if !storage.Exists(storageName) {
			fmt.Printf("Storage '%s' not found!\n", storageName)
			os.Exit(1)
		}

		var password string
		prompt := &survey.Password{
			Message: "Enter master password:",
		}
		survey.AskOne(prompt, &password)

		cfg, err := storage.LoadConfig(storageName)
		if err != nil {
			fmt.Printf("Failed to load storage: %v\n", err)
			os.Exit(1)
		}

		derivedKey := crypto.DeriveKey(password, cfg.Salt)
		if !isValidMasterPassword(derivedKey, cfg.Hash) {
			fmt.Println("Invalid master password!")
			os.Exit(1)
		}

		fmt.Println("Successfully authenticated!")
		//todo
	},
}

func isValidMasterPassword(derivedKey, storedHash []byte) bool {
	if len(derivedKey) != len(storedHash) {
		return false
	}

	return subtle.ConstantTimeCompare(derivedKey, storedHash) == 1
}

func init() {
	RootCmd.AddCommand(enterCmd)
}
