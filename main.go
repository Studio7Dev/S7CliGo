package main

import (
	"CLI/MCommands"
	misc "CLI/Misc"
	"CLI/s7cli/commands"
	"fmt"
	"os"
	"strings"
)

var (
	f_  = misc.Funcs{}
	MC_ = MCommands.MC{}
)

// import (
//
//	Auth "CLI/Auth"
//	HugginFace "CLI/HugAI"
//	Searx "CLI/SearXNG"
//	Movie_ "CLI/TMDB"
//	blackbox "CLI/blackbox"
//	MerlinAI "CLI/merlin_cli"
//	"CLI/s7cli/commands"
//
// )

func main() {

	settings, err := f_.LoadSettings()
	if err != nil {
		fmt.Println(err)
		return
	}
	if settings.Username == "" || settings.Password == "" {
		for {
			if settings.Username != "" && settings.Password != "" {
				break
			}
			data, err := f_.LoginForm()
			if err != nil {
				fmt.Println("Error getting login data:", err)
				return
			}
			settings.Username = data.Username
			settings.Password = data.Password
			err = f_.UpdateSettings(settings)
			if err != nil {
				fmt.Println("Error updating settings:", err)
				return
			}

		}
	}
	username := settings.Username
	password := settings.Password
	// Authenticate the user
	status, err := f_.LoginRequest(username, password)
	if err != nil {
		fmt.Println("Error logging in:", err)
		fmt.Println("Would you like to try again? (y/n)")
		yN := MC_.GetInput()
		if strings.ToLower(yN) == "y" {
			settings.Username = ""
			settings.Password = ""
			err = f_.UpdateSettings(settings)
			if err != nil {
				fmt.Println("Error updating settings:", err)
				return
			}
			main()
		} else {
			os.Exit(0)
		}
		return
	}
	if status != 200 {
		fmt.Println("Login failed with status code:", status)
		fmt.Println("Would you like to try again? (y/n)")
		yN := MC_.GetInput()
		if strings.ToLower(yN) == "y" {
			settings.Username = ""
			settings.Password = ""
			err = f_.UpdateSettings(settings)
			// handle err
			if err != nil {
				fmt.Println("Error updating settings:", err)
				return
			}
			main()
		} else {
			os.Exit(0)
		}
		return
	}
	if !f_.Authenticated() {
		fmt.Println("User is not authenticated. Please try again.")
		fmt.Println("Would you like to try again? (y/n)")
		yN := MC_.GetInput()
		if strings.ToLower(yN) == "y" {
			settings.Username = ""
			settings.Password = ""
			err = f_.UpdateSettings(settings)
			if err != nil {
				fmt.Println("Error updating settings:", err)
				return
			}
			main()
		} else {
			os.Exit(0)
		}
		return
	}
	//go func() {
	//	auth_status := f_.Authenticated()
	//	if !auth_status {
	//		fmt.Println("User is not authenticated. Exiting...")
	//		os.Exit(1)
	//	}
	//}()
	DefaultHandler := commands.DefaultHandler
	// User is authenticated, proceed with the application
	f_.Banner()

	DefaultHandler = MC_.Init(DefaultHandler)
	MC_.Run(DefaultHandler)
}
