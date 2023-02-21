#!/bin/zsh

export GOOS="windows"
export GOARCH="amd64"
go build -o ./versions/windows-amd64/ .

export GOOS="windows"
export GOARCH="arm64"
go build -o ./versions/windows-arm64/ .

export GOOS="linux"
export GOARCH="amd64"
go build -o ./versions/linux-amd64/ .

export GOOS="linux"
export GOARCH="arm64"
go build -o ./versions/linux-arm64/ .

export GOOS="darwin"
export GOARCH="amd64"
go build -o ./versions/macos-amd64/ .

export GOOS="darwin"
export GOARCH="arm64"
go build -o ./versions/macos-arm64/ .