package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

// loginCommand creates the "login" subcommand
func loginCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	fs := flag.NewFlagSet("testmyapp login", flag.ExitOnError)
	var (
		loginFlags struct {
			// Add flags specific to the "login" subcommand
			username string
		}
	)

	fs.StringVar(&loginFlags.username, "u", "", "u option for login command")

	return &ffcli.Command{
		Name:       "login",
		ShortUsage: "login [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			if loginFlags.username == "" {
				return errors.New("username is required")
			}
			fmt.Print("Enter Password: ")
			bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			password := string(bytePassword)
			fmt.Println() // Print a newline after the user presses enter
			return login(loginFlags.username, password, &c)
		},
	}
}

func login(username, password string, cfg *Config) error {
	c := NewCustomHTTPClient(apiHost, "", "")
	t, r, userID, err := c.Login(username, password)
	if err != nil {
		return err
	}
	err = cfg.UpdateTokens(username, t, r, userID)
	if err != nil {
		return err
	}
	fmt.Println("Login successful")
	return nil
}
