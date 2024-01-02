package clis

import (
	"os"
	"os/exec"

	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/moqsien/neobox/pkgs/run"
	"github.com/spf13/cobra"
)

func (that *Cli) neobox() {
	const (
		runnerCmd string = "runner"
		keeperCmd string = "keeper"
	)
	binPath, _ := os.Executable()
	// run box-server and box-keeper in daemon.
	box := vctrl.NewBox(
		exec.Command(binPath, vctrl.NeoBoxCmdName, runnerCmd),
		exec.Command(binPath, vctrl.NeoBoxCmdName, keeperCmd),
	)

	neoboxCmd := &cobra.Command{
		Use:     vctrl.NeoBoxCmdName,
		Aliases: []string{"nb"},
		Short:   "NeoBox related CLIs.",
		GroupID: that.groupID,
	}

	neoboxCmd.AddCommand(&cobra.Command{
		Use:     "shell",
		Aliases: []string{"s"},
		Short:   "Start NeoBox shell.",
		Run: func(cmd *cobra.Command, args []string) {
			box.StartShell()
		},
	})

	neoboxCmd.AddCommand(&cobra.Command{
		Use:     runnerCmd,
		Aliases: []string{"r"},
		Short:   "Start NeoBox server.",
		Run: func(cmd *cobra.Command, args []string) {
			box.StartClient() // server.
		},
	})

	neoboxCmd.AddCommand(&cobra.Command{
		Use:     keeperCmd,
		Aliases: []string{"k"},
		Short:   "Start NeoBox keeper.",
		Run: func(cmd *cobra.Command, args []string) {
			box.StartKeeper()
		},
	})

	autoStartCmd := &cobra.Command{
		Use:     vctrl.NeoBoxAutoStartCmdName,
		Aliases: []string{"as"},
		Short:   "This command is for auto-start on system booting.",
		Run: func(cmd *cobra.Command, args []string) {
			box.AutoStart(cmd, args...)
		},
	}
	autoStartCmd.Flags().BoolP(run.RestartUseDomain, "d", false, "Uses domains for edgetunnels.")
	autoStartCmd.Flags().BoolP(run.RestartForceSingbox, "s", false, "Force to use sing-box as client.")
	autoStartCmd.Flags().BoolP(run.RestartShowProxy, "p", false, "Shows proxy info.")
	autoStartCmd.Flags().BoolP(run.RestartShowConfig, "c", false, "Shows the config string.")
	neoboxCmd.AddCommand(autoStartCmd)

	genScriptCmd := &cobra.Command{
		Use:     "gen-script",
		Aliases: []string{"gs"},
		Short:   "Generates an autostart script for NeoBox.",
		Long:    "The generated script can be added to your system booting list.",
		Run: func(cmd *cobra.Command, args []string) {
			box.GenAutoStartScript()
		},
	}
	neoboxCmd.AddCommand(genScriptCmd)

	that.rootCmd.AddCommand(neoboxCmd)
}
