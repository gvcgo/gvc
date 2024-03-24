package main

import "github.com/gvcgo/gvc/cmd"

var (
	GitTag  string
	GitHash string
)

func main() {
	cli := cmd.NewCli(GitTag, GitHash)
	cli.Run()
}
