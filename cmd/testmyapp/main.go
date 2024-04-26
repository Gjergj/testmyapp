package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	apiHost             = "http://testmyapp.io:8080"
	configFileName      = "config.yaml"
	ModeForce      Mode = iota
)

const configDirName = "testmyapp.io"

func main() {
	ctx := context.Background()

	var (
		globalFlags = flag.NewFlagSet("cli", flag.ExitOnError)
	)

	// Define the root command
	root := &ffcli.Command{
		Name:       "testmyapp",
		ShortUsage: "testmyapp [flags] <subcommand>",
		FlagSet:    globalFlags,
		Subcommands: []*ffcli.Command{
			createProjectCommand(),
			loginCommand(),
			uploadCommand(),
			listProjectCommand(),
			watchCommand(),
			versionCommand(),
			deleteCommand(),
			signupCommand(),
			signoutCommand(),
		},
		Exec: run,
	}

	// Parse command-line arguments using ff
	err := root.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Run the root command
	if err := root.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		os.Exit(1)
	}
}
func run(ctx context.Context, args []string) error {
	return flag.ErrHelp
}
func ensureConfigFile() (string, error) {
	// Get the user's config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user's config directory:", err)
		return "", err
	}

	// Get the path to the config file
	configFilePath := filepath.Join(configDir, configDirName, configFileName)

	//check if config file exists
	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		// Create the config directory
		err = os.MkdirAll(filepath.Join(configDir, configDirName), 0755)
		if err != nil {
			fmt.Println("Error creating config directory:", err)
			return "", err
		}
		//crete config file
		_, err = os.Create(configFilePath)
		if err != nil {
			fmt.Println("Error creating config file:", err)
			return "", err
		}
	}
	return configFilePath, nil
}

func watchCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	fs := flag.NewFlagSet("cli watch", flag.ExitOnError)
	var userID string
	fs.StringVar(&userID, "u", "", "user name")

	return &ffcli.Command{
		Name:       "watch",
		ShortUsage: "watch [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			if c.IsEmpty() {
				fmt.Println("Please login first")
				return nil
			}
			createProject(userID, &c)
			watchFiles("./...", userID, &c)
			return nil
		},
	}
}

func uploadDir() string {
	var watchDir string
	watchDir = os.Getenv("UPLOAD_DIR")
	if watchDir == "" {
		var err error
		watchDir, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting working directory:", err)
		}
	}
	return watchDir
}

func getFiles(dir string) []string {

	//workaround for ./...
	if dir == "./..." {
		dir = "."
	}
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			i, err := d.Info()
			if err != nil {
				fmt.Printf("Warning: skipping file %s error getting file info: %s\n", d.Name(), err.Error())
				return nil
			}
			if i.Mode()&os.ModeSymlink != 0 {
				fmt.Printf("Warning: skipping file %s it's a symling\n", d.Name())
				return nil
			}
			if i.Size() == 0 {
				fmt.Printf("Warning: skipping file %s file size is 0\n", d.Name())
				return nil
			}
			path = strings.ReplaceAll(path, dir+string(filepath.Separator), "")
			files = append(files, path)
		}
		return nil
	})

	//entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return nil
	}
	return files
}
