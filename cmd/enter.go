package cmd

import (
	"crypto/subtle"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rivo/tview"
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

		runTUI()
	},
}

func runTUI() {
	app := tview.NewApplication()

	list := tview.NewList().
		AddItem("service", "login", '1', nil).
		AddItem("Entry 2", "", '2', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Entry 3", "", '3', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop() // Выход по нажатию 'q'
		})

	list.SetBorder(true).
		SetTitle(" Entries ").
		SetTitleAlign(tview.AlignLeft)

	details := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Select an entry to view details")

	details.SetBorder(true).
		SetTitle(" Details ").
		SetTitleAlign(tview.AlignLeft)

	// Обработчик выбора элемента в списке
	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		details.SetText(fmt.Sprintf(
			"[green]Entry: [white]%s\n[green]Description: [white]%s\n\n[yellow]Press 'q' to quit",
			mainText, secondaryText,
		))
	})

	flex := tview.NewFlex().
		AddItem(list, 0, 1, true).
		AddItem(details, 0, 2, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
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
