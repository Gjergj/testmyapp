package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"io"
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

	t, userID, userName := c.Token(userName)
	r, _ := c.RefreshToken(userName)
	deleted := false

	// if project name was not provided, assume we want to delete project in current directory
	if projectName == "" {
		pwd, _ := os.Getwd()
		for _, project := range c.Accounts[userName].Projects {
			if project.ProjectDir == pwd {
				c.RemoveProject(userName, project.ProjectName)
				projectName = project.ProjectName
				deleted = true
			}
		}
	}
	if deleted {
		err := c.Save()
		if err != nil {
			fmt.Println("Error saving config file:", err)
		}

		client := NewCustomHTTPClient(apiHost, t, r)
		serverURL := fmt.Sprintf("%s/v1/users/%s/projects/%s", apiHost, userID, projectName)
		response, err := client.Delete(serverURL, nil)

		defer response.Body.Close()
		if response.StatusCode != 200 {
			fmt.Println("Error deleting project:", response.Status)
			// Read the response body
			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return err
			}
			fmt.Println(string(responseBody))
			return nil
		}
	}
	return nil
}