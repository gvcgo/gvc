package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Cmder struct {
	*cli.App
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
	that.vneobox()
	that.vbrowser()
	that.vhomebrew()
	that.vhost()
	that.vgithub()
	that.vcloc()
	that.asciinema()

	that.vconf()
	that.version()
	that.showinfo()
	that.uninstall()
}
