package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"goPasswordManager/internal/crypto"
	"goPasswordManager/internal/storage"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new password storage",
	Run: func(cmd *cobra.Command, args []string) {
		var storageName string
		promptName := &survey.Input{
			Message: "Enter storage name:",
		}
		survey.AskOne(promptName, &storageName)

		if storage.Exists(storageName) {
			fmt.Printf("Storage '%s' already exists!\n", storageName)
			os.Exit(1)
		}

		qs := []*survey.Question{
			{
				Name:     "password",
				Prompt:   &survey.Password{Message: "Enter master password:"},
				Validate: survey.Required,
			},
			{
				Name:     "confirm",
				Prompt:   &survey.Password{Message: "Confirm master password:"},
				Validate: survey.Required,
			},
		}

		answers := struct {
			Password string
			Confirm  string
		}{}
		survey.Ask(qs, &answers)

		if answers.Password != answers.Confirm {
			fmt.Println("Passwords do not match!")
			os.Exit(1)
		}

		salt := crypto.GenerateSalt(32)
		hash := crypto.DeriveKey(answers.Password, salt)

		err := storage.Create(storageName, hash, salt)
		if err != nil {
			fmt.Printf("Failed to create storage: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Storage '%s' successfully created!\n", storageName)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
