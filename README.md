# Railway Deployment

This app is ready to deploy on Railway without Docker. Railway will automatically build the Go app and set the `PORT` environment variable. The server will listen on that port, as required by Railway.

No Dockerfile is needed. If you need to run locally, just use:

```
go run main.go
```

Or build and run:

```
go build -o go-api-test
./go-api-test
```

## Deploy Script Usage

To quickly add, commit, and push your changes, use the `deploy.cmd` script:

### In PowerShell (recommended on Windows 10/11):

```
.\deploy.cmd
```

### In Command Prompt (cmd.exe):

```
deploy
```

or

```
deploy.cmd
```

If you see a 'not recognized' error in PowerShell, always use the `./` or `./deploy.cmd` prefix to run scripts from the current directory.
