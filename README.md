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
