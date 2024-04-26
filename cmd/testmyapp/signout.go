package main

import (
	"context"
	"flag"
	"github.com/peterbourgon/ff/v3/ffcli"
	"log"
)

// signout Command creates the "signupCommand" subcommand
func signoutCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	fs := flag.NewFlagSet("testmyapp signout", flag.ExitOnError)
	var (
		loginFlags struct {
			// Add flags specific to the "login" subcommand
			username string
		}
	)

	fs.StringVar(&loginFlags.username, "u", "", "u option for login command")

	return &ffcli.Command{
		Name:       "signout",
		ShortUsage: "signout [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			c.Clear(loginFlags.username)
			return c.Save()
		},
	}
}
