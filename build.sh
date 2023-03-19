#!/bin/bash

# Build for Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o bin/app_64.exe main.go

# Build for macOS 64-bit
GOOS=darwin GOARCH=amd64 go build -o bin/app_64_darwin main.go

# Build for Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o bin/app_64_linux main.go
