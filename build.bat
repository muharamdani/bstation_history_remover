@echo off

rem Build for Windows 64-bit
set GOOS=windows
set GOARCH=amd64
go build -o app_64.exe main.go

rem Build for Mac 64-bit
set GOOS=darwin
set GOARCH=amd64
go build -o app_64 main.go

rem Build for Linux 64-bit
set GOOS=linux
set GOARCH=amd64
go build -o app_64 main.go
