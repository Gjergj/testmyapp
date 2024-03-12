## CLI for testmyapp.io

Still very much a work in progress.

## Getting started

#### Install on macOS
```bash
brew install gjergj/tap/testmyapp
```

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

#### Watch file changes as you develop
```bash
testmyapp watch
```
Refresh browser to see changes.

### Update
```bash
brew upgrade testmyapp
```