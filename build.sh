#!/bin/zsh

export GOOS="windows"
export GOARCH="amd64"
go build -o ./versions/windows-amd64/ .
echo "compilation for win-amd64 is done."

export GOOS="windows"
export GOARCH="arm64"
go build -o ./versions/windows-arm64/ .
echo "compilation for win-arm64 is done."

export GOOS="linux"
export GOARCH="amd64"
go build -o ./versions/linux-amd64/ .
echo "compilation for linux-amd64 is done."

export GOOS="linux"
export GOARCH="arm64"
go build -o ./versions/linux-arm64/ .
echo "compilation for linux-arm64 is done."

export GOOS="darwin"
export GOARCH="arm64"
go build -o ./versions/macos-arm64/ .
echo "compilation for darwin-arm64 is done."

export GOOS="darwin"
export GOARCH="amd64"
go build -o ./versions/macos-amd64/ .
echo "compilation for darwin-amd64 is done."