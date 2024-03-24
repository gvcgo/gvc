package cmd

import (
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/spf13/cobra"
)

const (
	GroupID string = "vm"
)

/*
CLIs
*/
type Cli struct {
	rootCmd *cobra.Command
	groupID string
	gitTag  string
	gitHash string
}

func NewCli(gitTag, gitHash string) (c *Cli) {
	c = &Cli{
		rootCmd: &cobra.Command{
			Short: "geek's valuable collections",
			Long:  "g <Command> <SubCommand> --flags args...",
		},
		groupID: GroupID,
		gitTag:  gitTag,
		gitHash: gitHash,
	}
	c.rootCmd.AddGroup(&cobra.Group{ID: c.groupID, Title: "Command list: "})
	c.initiate()
	return
}

func (c *Cli) initiate() {
	RegisterAsciinema(c)
}

func (that *Cli) Run() {
	if err := that.rootCmd.Execute(); err != nil {
		gprint.PrintError("%+v", err)
	}
}
