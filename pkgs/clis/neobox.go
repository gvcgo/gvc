package clis

import (
	"os"
	"os/exec"

	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func (that *Cli) neobox() {
	const (
		boxCmd    string = "neobox"
		runnerCmd string = "runner"
		keeperCmd string = "keeper"
	)
	binPath, _ := os.Executable()
	neobox := vctrl.NewBox(
		exec.Command(binPath, boxCmd, runnerCmd),
		exec.Command(binPath, boxCmd, keeperCmd),
	)

	neoboxCmd := &cobra.Command{
		Use:     boxCmd,
		Aliases: []string{"nb"},
		Short:   "NeoBox related CLIs.",
		GroupID: that.groupID,
	}

	neoboxCmd.AddCommand(&cobra.Command{
		Use:     "shell",
		Aliases: []string{"s"},
		Short:   "Start NeoBox shell.",
		Run: func(cmd *cobra.Command, args []string) {
			neobox.StartShell()
		},
	})

	neoboxCmd.AddCommand(&cobra.Command{
		Use:     runnerCmd,
		Aliases: []string{"r"},
		Short:   "Start NeoBox server.",
		Run: func(cmd *cobra.Command, args []string) {
			neobox.StartClient() // server.
		},
	})

	neoboxCmd.AddCommand(&cobra.Command{
		Use:     keeperCmd,
		Aliases: []string{"k"},
		Short:   "Start NeoBox keeper.",
		Run: func(cmd *cobra.Command, args []string) {
			neobox.StartKeeper()
		},
	})
	that.rootCmd.AddCommand(neoboxCmd)
}
