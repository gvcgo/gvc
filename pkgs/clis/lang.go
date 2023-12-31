package clis

import (
	"github.com/moqsien/gvc/pkgs/clis/langs"
	"github.com/spf13/cobra"
)

func (that *Cli) Register(cmd *cobra.Command) {
	cmd.GroupID = that.groupID
	that.rootCmd.AddCommand(cmd)
}

func (that *Cli) langs() {
	/*
		go
		protoc
		grpc
	*/
	langs.SetGo(that)
	/*
		java
		maven
		gradle
	*/
	langs.SetJava(that)

	/*
		python
		pyenv
		pyenv-win
	*/
	langs.SetPython(that)
}
