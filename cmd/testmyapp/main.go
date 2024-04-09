package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rjeczalik/notify"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	apiHost             = "http://testmyapp.io:8080"
	configFileName      = "config.yaml"
	ModeForce      Mode = iota
)

// TODO: On macOS, application-specific data and configuration files are typically stored in the ~/Library/Application Support/ directory, where ~ is the home directory of the current user.
// For example, if your application is named Manager, you would typically store its data in ~/Library/Application Support/Manager/.
// On Linux, application-specific data and configuration files are typically stored in the /etc/ directory for system-wide applications. For user-specific applications, data is usually stored in the user's home directory under a subdirectory that starts with a dot, for example ~/.myapp/.
// For Linux, it returns ~/.config. For macOS, it returns ~/Library/Application Support.
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
		},
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
		fmt.Println("Config file does not exist")
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
		fmt.Println("Created config file:", err)
	}
	return configFilePath, nil
}

func watchCommand() *ffcli.Command {
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	fs := flag.NewFlagSet("cli watch", flag.ExitOnError)
	var projectName, userID string
	fs.StringVar(&projectName, "p", "", "project name")
	fs.StringVar(&userID, "u", "", "user name")

	return &ffcli.Command{
		Name:       "watch",
		ShortUsage: "watch [flags]",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			watchFiles(uploadDir(uploadDirWatchCurrentRecursive), projectName, userID, &c)
			return nil
		},
	}
}

func watchFiles(dir, projectName, userID string, cfg *Config) {
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Create a channel to receive OS signals CTRL+C
	sigs := make(chan os.Signal, 1)

	// Set up a watchpoint listening for events within a directory tree rooted
	// at current working directory. Dispatch remove events to c.
	if err := notify.Watch(dir, c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	go func() {
		for {
			select {
			case ei := <-c:
				_ = ei
				//log.Println("Got event:", ei)
				time.Sleep(200 * time.Millisecond)
				files := getFiles(dir)
				uploadFiles(projectName, userID, files, cfg)
			}
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a CTRL+C signal is received
	<-sigs
}

type uploadDirType int

const (
	uploadAnyDirRecursive uploadDirType = iota
	uploadDirWatchCurrent
	uploadDirWatchCurrentRecursive
)

func uploadDir(uploadType uploadDirType) string {
	var watchDir string
	switch uploadType {
	case uploadAnyDirRecursive:
		watchDir = os.Getenv("UPLOAD_DIR")
		if watchDir == "" {
			var err error
			watchDir, err = os.Getwd()
			if err != nil {
				fmt.Println("Error getting working directory:", err)
			}
			//watchDir = filepath.Join(watchDir, "...")
		}
	case uploadDirWatchCurrent:
		watchDir = "."
	case uploadDirWatchCurrentRecursive:
		watchDir = "./..."
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
			path = strings.ReplaceAll(path, dir+"/", "")
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
