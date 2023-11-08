#!/bin/bash

echo "Build for Windows amd64..."
GOOS=windows GOARCH=amd64 go build -o ./release/template/qcs_v0.0.0_windows_amd64/server.exe server.go
GOOS=windows GOARCH=amd64 go build -o ./release/template/qcs_v0.0.0_windows_amd64/Init/init.exe ./Init/init.go
echo "Done."

echo "Build for Mac arm64..." 
GOOS=darwin GOARCH=arm64 go build -o ./release/template/qcs_v0.0.0_darwin_arm64/server server.go
GOOS=darwin GOARCH=arm64 go build -o ./release/template/qcs_v0.0.0_darwin_arm64/Init/init ./Init/init.go
echo "Done."

echo "Build for Mac amd64"
GOOS=darwin GOARCH=amd64 go build -o ./release/template/qcs_v0.0.0_darwin_amd64/server server.go  
GOOS=darwin GOARCH=amd64 go build -o ./release/template/qcs_v0.0.0_darwin_amd64/Init/init ./Init/init.go
echo "Done."

echo "Build for Linux amd64"
GOOS=linux GOARCH=amd64 go build -o ./release/template/qcs_v0.0.0_linux_amd64/server server.go
GOOS=linux GOARCH=amd64 go build -o ./release/template/qcs_v0.0.0_linux_amd64/Init/init ./Init/init.go
echo "Done." 

echo "Build for Linux arm64 v8"
GOOS=linux GOARCH=arm64 go build -o ./release/template/qcs_v0.0.0_linux_arm64/server server.go
GOOS=linux GOARCH=arm64 go build -o ./release/template/qcs_v0.0.0_linux_arm64/Init/init ./Init/init.go
echo "Done."