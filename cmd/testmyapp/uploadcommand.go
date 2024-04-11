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
	"path/filepath"
)

func uploadCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	fs := flag.NewFlagSet("testmyapp upload", flag.ExitOnError)

	var userName string
	fs.StringVar(&userName, "u", "", "user name")

	return &ffcli.Command{
		Name:       "upload",
		ShortUsage: "upload [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			if c.IsEmpty() {
				fmt.Println("Please login first")
				return nil
			}
			createProject(userName, &c)
			files := getFiles(uploadDir())
			uploadFiles(userName, files, &c)
			return nil
		},
	}
}

func uploadFiles(userName string, files []string, c *Config) {
	if len(files) > models.MaxUploadFiles {
		fmt.Println(fmt.Sprintf("Maximum number of files allowed is %d", models.MaxUploadFiles))
		return
	}

	totalSize := 0
	foundIndex := false
	filesToUpload := make([]string, 0)
	for i, file := range files {
		ext := filepath.Ext(file)
		if !models.AllowedFileType(ext) {
			fmt.Println(fmt.Sprintf("File type %s is not allowed on file %s. Will not be uploaded.", ext, file))
			continue
		}

		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Println("Error getting file info:", err)
			return
		}
		fileName := fileInfo.Name()
		if fileInfo.Size() > models.MaxFileSizeLimit {
			fmt.Println(fmt.Sprintf("%s file size %d exceeds the limit of %d. Will not be uploaded", fileName, fileInfo.Size(), models.MaxFileSizeLimit))
			continue
		}
		totalSize += int(fileInfo.Size())
		// check if the file name is too long
		if len(fileName) > models.MaxFileNameLength {
			fmt.Println(fmt.Sprintf("File name %s is too long. Allowed file name length is %d. Will not be uploaded.", fileName, models.MaxFileNameLength))
			continue
		}
		if fileInfo.Name() == "index.html" {
			foundIndex = true
		}
		filesToUpload = append(filesToUpload, files[i])
	}

	if !foundIndex {
		fmt.Println("This directory does not contain index.html. An index.html file is required")
		return
	}
	if totalSize > models.MaxUploadSize {
		fmt.Println(fmt.Sprintf("Total size of files %d exceeds the limit of %d", totalSize, models.MaxUploadSize))
		return
	}

	fmt.Println("Uploading files...", filesToUpload)

	t, userID, userName := c.Token(userName)
	r, _ := c.RefreshToken(userName)

	if userName == "" || userID == "" || t == "" || r == "" {
		fmt.Println("Please login first")
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var ok bool
	var projectName string
	if projectName, ok = c.GetProjectID(userID, dir); !ok {
		fmt.Println("could not find project for current directory")
		return
	}

	serverURL := fmt.Sprintf("%s/v1/users/%s/projects/%s", apiHost, userID, projectName)

	cl := NewCustomHTTPClient(apiHost, t, r)

	// Make the POST request
	response, err := cl.Upload(serverURL, filesToUpload)
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
		err = c.UpdateTokens(userName, cl.Token, cl.RefreshToken, userID)
		if err != nil {
			return
		}
	}
	// Print the response
	fmt.Println(apiResp.Message)
	fmt.Println()
	fmt.Println("Refresh the page to see the changes.")
	fmt.Println("Try 'testmyapp watch' it will detect changes and automatically upload files.")
}
