package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"goPasswordManager/internal/crypto"
	"goPasswordManager/internal/models"
	"goPasswordManager/internal/storage"
)

var (
	version     string
	description string
)

var addCmd = &cobra.Command{
	Use:   "add [name] [service] [login]",
	Short: "Add a new entry to the storage",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		service := args[1]
		login := args[2]

		if !storage.Exists(name) {
			fmt.Printf("Storage '%s' not found!\n", name)
			os.Exit(1)
		}

		var masterPassword string
		prompt := &survey.Password{
			Message: "Enter master masterPassword:",
		}
		survey.AskOne(prompt, &masterPassword)

		cfg, err := storage.LoadConfig(name)
		if err != nil {
			fmt.Printf("Failed to load storage: %v\n", err)
			os.Exit(1)
		}

		derivedKey := crypto.DeriveKey(masterPassword, cfg.Salt)
		if !isValidMasterPassword(derivedKey, cfg.Hash) {
			fmt.Println("Invalid master masterPassword!")
			os.Exit(1)
		}

		password := models.Password{
			Service:     service,
			Login:       login,
			Version:     version,
			Description: description,
		}

		if err := storage.AddEntry(name, password); err != nil {
			fmt.Printf("Failed to add password: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Entry added successfully!")
	},
}

func init() {
	addCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the entry")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the entry")

	RootCmd.AddCommand(addCmd)
}
