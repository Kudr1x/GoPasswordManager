package cmd

import (
	"crypto/subtle"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"goPasswordManager/internal/auth"
	"goPasswordManager/internal/models"
	"os"

	"github.com/atotto/clipboard"
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

		var masterPassword string
		prompt := &survey.Password{
			Message: "Enter master masterPassword:",
		}
		survey.AskOne(prompt, &masterPassword)

		cfg, err := storage.LoadConfig(storageName)
		if err != nil {
			fmt.Printf("Failed to load storage: %v\n", err)
			os.Exit(1)
		}

		derivedKey := crypto.DeriveKey(masterPassword, cfg.Salt)
		if !isValidMasterPassword(derivedKey, cfg.Hash) {
			fmt.Println("Invalid master masterPassword!")
			os.Exit(1)
		}

		runTUI(cfg.Entries, masterPassword)
		auth.ZeroString(masterPassword)
	},
}

func runTUI(entries []models.Password, masterPass string) {
	app := tview.NewApplication()

	auth.StorePassword(masterPass)
	defer auth.GetPassword()

	flex := tview.NewFlex()

	list := tview.NewList()

	for i, entry := range entries {
		mainText := fmt.Sprintf("[%d] %s", i+1, entry.Service)
		secondaryText := fmt.Sprintf("Login: %s", entry.Login)
		list.AddItem(mainText, secondaryText, rune('1'+i), nil)
	}

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index >= len(entries) {
			return
		}

		generatedPass := generatePassword(entries[index])

		if err := clipboard.WriteAll(generatedPass); err != nil {
			showMessage(app, err.Error(), flex)
		} else {
			showMessage(app, "Password copied to clipboard!", flex)
		}
	})

	list.AddItem("Exit", "Click to close the application", 'q', func() {
		app.Stop()
	})

	list.SetBorder(true).
		SetTitle(" Passwords ").
		SetTitleAlign(tview.AlignLeft)

	details := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Select password to view")

	details.SetBorder(true).
		SetTitle(" Info ").
		SetTitleAlign(tview.AlignLeft)

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index >= len(entries) {
			return
		}

		entry := entries[index]
		detailsText := fmt.Sprintf(
			"[green]Service: [white]%s\n"+
				"[green]Login: [white]%s\n"+
				"[green]Version: [white]%s\n"+
				"[green]Description: [white]%s\n\n"+
				"[yellow]Press Enter to copy\n"+
				"[yellow]Press Q to exit",
			entry.Service,
			entry.Login,
			entry.Version,
			entry.Description,
		)
		details.SetText(detailsText)
	})

	flex.
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

func generatePassword(password models.Password) string {
	//todo
	return "ugabuga"
}

func showMessage(app *tview.Application, text string, flex *tview.Flex) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			app.SetRoot(flex, true)
		})
	app.SetRoot(modal, false)
}

func init() {
	RootCmd.AddCommand(enterCmd)
}
