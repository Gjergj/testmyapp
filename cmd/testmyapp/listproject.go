package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Gjergj/testmyapp/pkg/models"
	"github.com/peterbourgon/ff/v3/ffcli"
	"io"
)

func listProjectCommand(c *Config) *ffcli.Command {
	fs := flag.NewFlagSet("cli list", flag.ExitOnError)
	//var (
	//	createFlags struct {
	//		// Add flags specific to the "create" subcommand
	//		foo string
	//	}
	//)

	var username string

	//fs.StringVar(&createFlags.foo, "foo", "", "Foo option for create command")
	//fs.StringVar(&p.ProjectName, "p", "", "project name")
	fs.StringVar(&username, "u", "", "user name")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "create [flags]",
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
	//fmt.Println("Token:", t)
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

	// Check the response status
	//if response.StatusCode != http.StatusOK {
	//	fmt.Printf("HTTP request failed with status code: %d\n", response.StatusCode)
	//	return
	//}

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

	for _, p := range apiResp.Projects {
		fmt.Printf("%s\t%s\n", p.ProjectName, p.URL)
	}
}