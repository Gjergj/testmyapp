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
func createProjectCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	fs := flag.NewFlagSet("testmyapp create", flag.ExitOnError)

	username := ""

	fs.StringVar(&username, "u", "", "user name")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "create [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			createProject(username, &c)
			return nil
		},
	}
}

func createProject(username string, c *Config) {
	t, userID, username := c.Token(username)
	r, _ := c.RefreshToken(username)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	if username == "" || userID == "" || t == "" || r == "" {
		fmt.Println("Please login first")
		return
	}

	// check if there's a project for the user in the current directory
	if _, ok := c.GetProjectID(userID, path); ok {
		fmt.Println("Project already exists in this directory")
		fmt.Println("Try running the 'testmyapp list' command to see all your projects")
		fmt.Println("Run 'testmyapp delete' to delete the project in this directory")
		return
	}

	client := NewCustomHTTPClient(apiHost, t, r)
	serverURL := fmt.Sprintf("%s/v1/users/%s/projects", apiHost, userID)

	response, err := client.Post(serverURL, nil)
	if err != nil {
		fmt.Println("Error creating project:", err)
		return
	}
	defer response.Body.Close()
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
