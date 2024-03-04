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

func uploadCommand(c *Config) *ffcli.Command {
	fs := flag.NewFlagSet("testmyapp upload", flag.ExitOnError)

	var projectName, userName string
	fs.StringVar(&projectName, "p", "", "project name")
	fs.StringVar(&userName, "u", "", "user name")

	return &ffcli.Command{
		Name:       "upload",
		ShortUsage: "upload [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			files := getFiles(uploadDir(uploadAnyDirRecursive))
			uploadFiles(projectName, userName, files, c)
			return nil
		},
	}
}

func uploadFiles(projectName, userName string, files []string, c *Config) {
	fmt.Println("Uploading files...", files)

	t, userID, userName := c.Token(userName)
	r, _ := c.RefreshToken(userName)

	if projectName == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		var ok bool
		if projectName, ok = c.GetProjectID(userID, dir); !ok {
			fmt.Println("could not find project for current directory")
			return
		}
	}

	serverURL := fmt.Sprintf("%s/v1/users/%s/projects/%s", apiHost, userID, projectName)

	cl := NewCustomHTTPClient(apiHost, t, r)

	// Make the POST request
	response, err := cl.Upload(serverURL, files)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	apiResp := models.UploadFilesResponse{}
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	// Check if the token has changed
	if cl.Token != t {
		c.UpdateTokens(userName, cl.Token, cl.RefreshToken, userID)
	}
	// Print the response
	fmt.Println(apiResp.Message)
}
