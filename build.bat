@echo off

echo Build for Windows amd64...
set GOOS=windows
set GOARCH=amd64
go build -o ./release/template/qcs_v0.0.0_windows_amd64/server.exe server.go
go build -o ./release/template/qcs_v0.0.0_windows_amd64/Init/init.exe ./Init/init.go 
echo Done.

echo Build for Mac arm64...
set GOOS=darwin
set GOARCH=arm64
go build -o ./release/template/qcs_v0.0.0_darwin_arm64/server server.go
go build -o ./release/template/qcs_v0.0.0_darwin_arm64/Init/init ./Init/init.go
echo Done.

echo Build for Mac amd64
set GOOS=darwin
set GOARCH=amd64
go build -o ./release/template/qcs_v0.0.0_darwin_amd64/server server.go
go build -o ./release/template/qcs_v0.0.0_darwin_amd64/Init/init ./Init/init.go
echo Done.

echo Build for Linux amd64
set GOOS=linux
set GOARCH=amd64
go build -o ./release/template/qcs_v0.0.0_linux_amd64/server server.go
go build -o ./release/template/qcs_v0.0.0_linux_amd64/Init/init ./Init/init.go
echo Done.

echo Build for Linux arm64 v8
set GOOS=linux
set GOARCH=arm64
go build -o ./release/template/qcs_v0.0.0_linux_arm64/server server.go
go build -o ./release/template/qcs_v0.0.0_linux_arm64/Init/init ./Init/init.go
echo Done.