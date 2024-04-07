package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// getConfig gets the config file from the user's config directory
func getConfig() (Config, error) {
	f, err := ensureConfigFile()
	if err != nil {
		return Config{}, err
	}
	// Open the config file
	configFile, err := os.Open(f)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return Config{}, err
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Println("Error closing config file:", err)
		}
	}(configFile)
	cfg := Config{
		Accounts: make(map[string]Account),
	}
	err = yaml.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return Config{}, err
	}
	return cfg, nil
}

// Save saves the config file to the user's config directory
func (c *Config) Save() error {
	f, err := ensureConfigFile()
	if err != nil {
		return err
	}
	// Open the config file
	configFile, err := os.OpenFile(f, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return err
	}
	defer configFile.Close()
	err = yaml.NewEncoder(configFile).Encode(c)
	if err != nil {
		fmt.Println("Error encoding config file:", err)
	}
	return err
}

type Config struct {
	Accounts map[string]Account `yaml:"accounts"`
}

type Account struct {
	Token        string    `yaml:"token"`
	RefreshToken string    `yaml:"refresh_token"`
	Projects     []Project `yaml:"projects"`
	UserID       string    `yaml:"user_id"`
}

func (a *Account) ToString() string {
	return fmt.Sprintf("Token: %s, Projects: %v", a.Token, a.Projects)
}

func (c *Config) RemoveProject(username string, projectName string) {
	a, ok := c.Accounts[username]
	if !ok {
		fmt.Printf("Account %s does not exist", username)
		return
	}
	a.RemoveProject(projectName)
	c.Accounts[username] = a
}

func (c *Config) UpdateTokens(username string, t string, r string, userID string) error {
	a, ok := c.Accounts[username]
	if ok {
		a.Token = t
		a.RefreshToken = r
		a.UserID = userID
	} else {
		a = Account{
			UserID:       userID,
			Token:        t,
			RefreshToken: r,
		}
	}
	c.Accounts[username] = a
	err := c.Save()
	if err != nil {
		return fmt.Errorf("Error saving config file %w", err)
	}
	return nil
}

func (c *Config) GetProjectID(userID, directory string) (string, bool) {
	for _, account := range c.Accounts {
		if account.UserID == userID {
			return account.DirectoryProject(directory)
		}
	}
	return "", false
}

func (c *Config) Token(username string) (string, string, string) {
	//username not specified and more than one account exists
	if len(c.Accounts) > 1 && username == "" {
		fmt.Println("Please specify an account")
		return "", "", ""
	}

	p, ok := c.Accounts[username]
	if !ok {
		//account not found and more than one account exists
		if len(c.Accounts) == 0 {
			fmt.Println("No accounts exist, please login")
			return "", "", ""
		}
		// return the first account
		for k, v := range c.Accounts {
			return v.Token, v.UserID, k
		}
	}
	return p.Token, p.UserID, username
}

func (c *Config) RefreshToken(username string) (string, string) {
	//username not specified and more than one account exists
	if len(c.Accounts) > 1 && username == "" {
		fmt.Println("Please specify an account")
		return "", ""
	}

	p, ok := c.Accounts[username]
	if !ok {
		//account not found and more than one account exists
		if len(c.Accounts) == 0 {
			fmt.Println("No accounts exist, please login")
			return "", ""
		}
		// return the first account
		for _, v := range c.Accounts {
			return v.RefreshToken, v.UserID
		}
	}
	return p.RefreshToken, p.UserID
}

func (c *Config) AddProject(accountName string, project Project, mode Mode) {
	//c.Accounts[accountName].AddProject(project, mode)
	p, ok := c.Accounts[accountName]
	if !ok {
		fmt.Printf("Account %s does not exist", accountName)
		return
	}
	p.AddProject(project, mode)
	c.Accounts[accountName] = p
}

type Mode int

func (a *Account) AddProject(project Project, mode Mode) {
	if mode == ModeForce {
		a.RemoveProject(project.ProjectName)
	} else if p, ok := a.DirectoryHasProject(project.ProjectDir); !ok {
		fmt.Printf("Directory %s already exists for project %s", project.ProjectDir, p)
		return
	}
	a.Projects = append(a.Projects, project)
}

func (a *Account) RemoveProject(projectName string) {
	for i, project := range a.Projects {
		if project.ProjectName == projectName {
			a.Projects = append(a.Projects[:i], a.Projects[i+1:]...)
			return
		}
	}
}

func (a *Account) ProjectDirectory(projectName string) (string, bool) {
	for _, project := range a.Projects {
		if project.ProjectName == projectName {
			return project.ProjectDir, true
		}
	}
	return "", false
}

func (a *Account) DirectoryProject(dir string) (string, bool) {
	for _, project := range a.Projects {
		if project.ProjectDir == dir {
			return project.ProjectName, true
		}
	}
	return "", false
}

func (a *Account) DirectoryHasProject(dirName string) (string, bool) {
	for _, p := range a.Projects {
		if p.ProjectDir == dirName {
			return p.ProjectName, true
		}
	}
	return "", false
}

type Project struct {
	//mapping of project name to directory
	ProjectName string `yaml:"project_name"`
	ProjectDir  string `yaml:"project_dir"`
}

func (p *Project) ToString() string {
	return fmt.Sprintf("ProjectName: %s, ProjectDir: %s", p.ProjectName, p.ProjectDir)
}
