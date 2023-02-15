package main

import (
	"os"

	"github.com/moqsien/gvc/pkgs/cmd"
)

func main() {
	c := cmd.New()
	c.Run(os.Args)
}
