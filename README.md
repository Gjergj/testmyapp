## CLI for testmyapp.io

Still very much a work in progress.

## Getting started

### Installation
#### macOS
```bash
brew install gjergj/tap/testmyapp
```

#### Windows
```bash
winget install testmyapp
```
#### Linux
Locate the latest release from the [releases page](https://github.com/Gjergj/testmyapp/releases) and download the appropriate package for your system.

Ubuntu example:
```bash
wget https://github.com/Gjergj/testmyapp/releases/download/v0.0.68/testmyapp_0.0.68_amd64.deb
sudo dpkg -i testmyapp_0.0.68_amd64.deb
```


### Commands
#### Login
```bash
testmyapp login -u=<username>
```

#### Upload your web site
Create an `index.html` file in the current directory and upload it to the project.
```bash
testmyapp upload
```

#### List your projects
```bash
testmyapp list
```

#### Delete a project.
Will only delete the project from the testmyapp.io server. It will not delete anything machine.
```bash
testmyapp delete
```
Will delete a project that is in the current directory.
To delete a specific project:
```bash
testmyapp delete -p=<project-id>
```

#### Watch file changes as you develop
This will watch for changes in the current directory and upload the changes to the server.
```bash
testmyapp watch
```
Refresh browser to see changes.

### Update
```bash
brew update
brew upgrade testmyapp
```