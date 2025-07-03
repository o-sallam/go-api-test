# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o manage-system-api

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/manage-system-api .
COPY --from=builder /app/wwwroot ./wwwroot
EXPOSE 8080
CMD ["./manage-system-api"]
