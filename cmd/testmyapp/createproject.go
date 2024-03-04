package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Gjergj/testmyapp/pkg/models"
	"github.com/peterbourgon/ff/v3/ffcli"
	"io"
	"log"
	"os"
)

// createProjectCommand creates a new project for user
func createProjectCommand(c *Config) *ffcli.Command {
	fs := flag.NewFlagSet("cli create", flag.ExitOnError)

	username := ""

	fs.StringVar(&username, "u", "", "user name")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "create [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			createProject(username, c)
			return nil
		},
	}
}

func createProject(username string, c *Config) {
	t, userID, username := c.Token(username)
	r, _ := c.RefreshToken(username)
	//fmt.Println("Token:", t)
	client := NewCustomHTTPClient(apiHost, t, r)
	serverURL := fmt.Sprintf("%s/v1/users/%s/projects", apiHost, userID)

	response, err := client.Post(serverURL, nil)

	defer response.Body.Close()
	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	// Parse the response body to get the JWT and refresh token
	apiResp := models.CreateProjectResponse{}
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	if pID, ok := c.GetProjectID(userID, path); ok {
		c.RemoveProject(username, pID)
	}

	c.AddProject(username, Project{
		ProjectName: apiResp.ProjectName,
		ProjectDir:  path,
	}, ModeForce)

	err = c.Save()
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}

	// Print the response
	fmt.Println("New Project:", apiResp.Message)
}
