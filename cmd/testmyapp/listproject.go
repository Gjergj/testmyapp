package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Gjergj/testmyapp/pkg/models"
	"github.com/peterbourgon/ff/v3/ffcli"
	"io"
	"os"
)

func listProjectCommand(c *Config) *ffcli.Command {
	fs := flag.NewFlagSet("testmyapp list", flag.ExitOnError)

	var username string
	fs.StringVar(&username, "u", "", "user name")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "list [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			getAllProjectsByUserID(username, c)
			return nil
		},
	}
}

// method that makes a get request to the server to fetch all projects for a user where user id is in the header
func getAllProjectsByUserID(username string, c *Config) {
	t, userID, userName := c.Token(username)
	r, _ := c.RefreshToken(username)
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
		c.UpdateTokens(userName, cl.Token, cl.RefreshToken, userID)
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	for _, project := range apiResp.Projects {
		fmt.Printf(project.URL)
		found := false
		for _, p := range c.Accounts[userName].Projects {
			if p.ProjectName == project.ProjectName && p.ProjectDir == pwd {
				// Print the current directory with an arrow
				fmt.Printf("\t←")
				found = true
			} else if p.ProjectName == project.ProjectName {
				// exists in this account but not in this pc
				found = true
			}
		}
		if !found {
			// not found in this account
			fmt.Printf("\t❌")
		}
		fmt.Println()
	}
}
