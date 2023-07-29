package cmd

import (
	"os"
	"os/exec"

	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

func (that *Cmder) vneobox() {
	const (
		neoRunner string = "neobox-runner"
		neoKeeper string = "neobox-keeper"
	)
	binPath, _ := os.Executable()
	neobox := vctrl.NewBox(exec.Command(binPath, neoRunner), exec.Command(binPath, neoKeeper))

	commands := &cli.Command{
		Name:    "neobox-shell",
		Aliases: []string{"shell", "box", "ns"},
		Usage:   "Start a neobox shell.",
		Action: func(ctx *cli.Context) error {
			neobox.StartShell()
			return nil
		},
	}
	that.Commands = append(that.Commands, commands)

	commandr := &cli.Command{
		Name:    neoRunner,
		Aliases: []string{"nbrunner", "nbr"},
		Usage:   "Start a neobox client.",
		Action: func(ctx *cli.Context) error {
			neobox.StartClient()
			return nil
		},
	}
	that.Commands = append(that.Commands, commandr)

	commandk := &cli.Command{
		Name:    neoKeeper,
		Aliases: []string{"nbkeeper", "nbk"},
		Usage:   "Start a neobox keeper.",
		Action: func(ctx *cli.Context) error {
			neobox.StartKeeper()
			return nil
		},
	}
	that.Commands = append(that.Commands, commandk)
}
