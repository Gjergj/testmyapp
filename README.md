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

### Commands
#### Login
```bash
testmyapp login -u=<username>
```

#### Create new project
This will create a new project and return the project URL. When visiting the URL you will see a 'Hello World' website.
```bash
testmyapp create
```
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