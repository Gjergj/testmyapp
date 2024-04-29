package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Gjergj/testmyapp/pkg/models"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/peterbourgon/ff/v3/ffcli"
	"io"
	"log"
	"os"
)

const (
	existsOnlyLocally  exists = "Local"
	existsOnlyOnRemote exists = "Remote"
	existsOnBoth       exists = "Both"
	existsHere         exists = "Current Directory"
)

type exists string

func listProjectCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	fs := flag.NewFlagSet("testmyapp list", flag.ExitOnError)

	var username string
	var printDir bool
	fs.StringVar(&username, "u", "", "user name")
	fs.BoolVar(&printDir, "d", false, "print directories")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "list [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			if c.IsEmpty() {
				fmt.Println("Please login first")
				return nil
			}
			getAllProjectsByUserID(username, printDir, &c)
			return nil
		},
	}
}

// method that makes a get request to the server to fetch all projects for a user where user id is in the header
func getAllProjectsByUserID(username string, printDirs bool, c *Config) {
	t, userID, userName := c.Token(username)
	r, _ := c.RefreshToken(username)
	if userName == "" || userID == "" || t == "" || r == "" {
		fmt.Println("Please login first")
		return
	}
	cl := NewCustomHTTPClient(apiHost, t, r)

	// URL to send the GET request to
	serverURL := fmt.Sprintf("%s/v1/users/%s/projects", apiHost, userID)
	response, err := cl.Get(serverURL)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return
	}

	defer response.Body.Close()
	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	// Parse the response body to get the JWT and refresh token
	apiResp := models.GetProjectsResponse{}
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	if apiResp.Status != 0 {
		fmt.Println("Error:", apiResp.Message)
		return
	}
	// Check if the token has changed
	if cl.Token != t {
		err = c.UpdateTokens(userName, cl.Token, cl.RefreshToken, userID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	type p struct {
		ProjectURL  string
		ProjectDir  string
		ProjectName string
		Exists      exists
	}
	projects := make(map[string]p)

	for _, project := range apiResp.Projects {
		projects[project.ProjectName] = p{
			ProjectURL:  project.URL,
			ProjectName: project.ProjectName,
			ProjectDir:  "",
			Exists:      existsOnlyOnRemote,
		}
		for _, rp := range c.Accounts[userName].Projects {
			if rp.ProjectName == project.ProjectName && rp.ProjectDir == pwd {
				projects[project.ProjectName] = p{
					ProjectURL:  project.URL,
					ProjectName: project.ProjectName,
					ProjectDir:  rp.ProjectDir,
					Exists:      existsHere,
				}
			} else if rp.ProjectName == project.ProjectName {
				projects[project.ProjectName] = p{
					ProjectURL:  project.URL,
					ProjectName: project.ProjectName,
					ProjectDir:  rp.ProjectDir,
					Exists:      existsOnBoth,
				}
			}
		}
	}
	for _, rp := range c.Accounts[userName].Projects {
		if _, ok := projects[rp.ProjectName]; !ok {
			projects[rp.ProjectName] = p{
				ProjectURL:  "",
				ProjectName: rp.ProjectName,
				ProjectDir:  rp.ProjectDir,
				Exists:      existsOnlyLocally,
			}
		}
	}
	projectTable := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))).Width(120)
	if printDirs {
		projectTable.Headers("Exists", "Project Name", "URL", "Dir")
	} else {
		projectTable.Headers("Exists", "Project Name", "URL")
	}

	// Print the projects
	for _, project := range projects {
		if printDirs {
			projectTable.Row(string(project.Exists), project.ProjectName, project.ProjectURL, project.ProjectDir)
		} else {
			projectTable.Row(string(project.Exists), project.ProjectName, project.ProjectURL)
		}
	}
	fmt.Println(projectTable)
}
