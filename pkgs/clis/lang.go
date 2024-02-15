package clis

import (
	"github.com/gvcgo/gvc/pkgs/clis/langs"
	"github.com/spf13/cobra"
)

// register new command to Cli.
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

	/*
		NodeJS
	*/
	langs.SetNodeJS(that)

	/*
		Flutter
		Flutter-Android-DEV using VSCode.
	*/
	langs.SetFlutter(that)

	/*
		julia
	*/
	langs.SetJulia(that)

	/*
		rust
		rust acceleration.
	*/
	langs.SetRust(that)

	/*
		cygwin
		msys2
		vcpkg
	*/
	langs.SetCpp(that)

	/*
		typst
	*/
	langs.SetTypst(that)

	/*
		vlang
		v-analyzer
		v-analyzer extension for VSCode
	*/
	langs.SetVlang(that)

	/*
		zig
		zls
	*/
	langs.SetZig(that)
}
