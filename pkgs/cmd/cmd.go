package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Cmder struct {
	*cli.App
	gitTag  string
	gitHash string
	gitTime string
}

func New() *Cmder {
	c := &Cmder{
		App: &cli.App{
			Usage:       "gvc <Command> <SubCommand>...",
			Description: "A productive tool to manage your development environment.",
			Commands:    []*cli.Command{},
		},
	}
	c.initiate()
	return c
}

func (c *Cmder) RunApp() {
	args := HandleArgs(os.Args...)
	c.Run(args)
}

func (that *Cmder) initiate() {
	that.vgo()
	that.vprotobuf()
	that.vpython()
	that.vjava()
	that.vmaven()
	that.vgradle()
	that.vnodejs()
	that.vflutter()
	that.vjulia()
	that.vrust()
	that.vcpp()
	that.vtypst()
	that.vlang()

	that.vscode()
	that.vnvim()
	that.vgpt()
	that.vneobox()
	that.vbrowser()
	that.vhomebrew()
	that.vgsudo()
	that.vhost()
	that.vgit()
	that.vinstallGitWin()
	that.vgithub()
	that.vcloc()
	that.vasciinema()
	that.vdocker()

	that.vconf()
	that.vsshFiles()
	that.version()
	that.checkUpdate()
	that.showinfo()
	that.uninstall()
}
