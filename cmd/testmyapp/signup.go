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

// signupCommand creates the "signupCommand" subcommand
func signupCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	fs := flag.NewFlagSet("testmyapp signup", flag.ExitOnError)
	var (
		loginFlags struct {
			// Add flags specific to the "login" subcommand
			username string
		}
	)

	fs.StringVar(&loginFlags.username, "u", "", "u option for login command")

	return &ffcli.Command{
		Name:       "signup",
		ShortUsage: "signup [flags]",
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
			return signup(loginFlags.username, password, &c)
		},
	}
}

func signup(username, password string, cfg *Config) error {
	c := NewCustomHTTPClient(apiHost, "", "")
	userID, err := c.SignUp(username, password)
	if err != nil {
		fmt.Println("Error signing up:", err)
		return err
	}
	fmt.Println("You should have received an email with a verification code.")
	fmt.Println("Enter the code below to verify your email address.")
	fmt.Print("Verification Code: ")
	var code string
	_, err = fmt.Scanln(&code)
	if err != nil {
		return err
	}
	t, r, userID, err := c.VerifyOTP(userID, code)
	if err != nil {
		return err
	}
	err = cfg.UpdateTokens(username, t, r, userID)
	if err != nil {
		return err
	}
	fmt.Println("Signup successful")
	fmt.Println()
	fmt.Println("You are also signed in.")
	fmt.Println("Now go to a directory with your static site in it and run 'testmyapp upload' to get a URL to view your site.")
	return nil
}
