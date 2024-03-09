package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"os"
)

func deleteCommand(c *Config) *ffcli.Command {
	fs := flag.NewFlagSet("testmyapp delete", flag.ExitOnError)

	var projectName, userName string
	fs.StringVar(&projectName, "p", "", "project name")
	fs.StringVar(&userName, "u", "", "user name")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "delete [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			err := deleteProject(projectName, userName, c)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func deleteProject(projectName, userName string, c *Config) error {

	_, _, userName = c.Token(userName)
	//r, _ := c.RefreshToken(userName)

	if projectName == "" {
		pwd, _ := os.Getwd()
		for _, project := range c.Accounts[userName].Projects {
			if project.ProjectDir == pwd {
				c.RemoveProject(userName, project.ProjectName)
			}
		}
	}
	err := c.Save()
	if err != nil {
		fmt.Println("Error saving config file:", err)
	}
	return err
}
