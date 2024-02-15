package langs

import (
	"github.com/gvcgo/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func SetNodeJS(reg IRegister) {
	nodeCmd := &cobra.Command{
		Use:     "nodejs",
		Aliases: []string{"n"},
		Short:   "NodeJS related CLIs.",
	}

	remoteCmd := &cobra.Command{
		Use:     "remote",
		Aliases: []string{"r"},
		Short:   "Shows available versions from remote website.",
		Run: func(cmd *cobra.Command, args []string) {
			nv := vctrl.NewNodeVersion()
			nv.ShowVersions()
		},
	}
	nodeCmd.AddCommand(remoteCmd)

	useCmd := &cobra.Command{
		Use:     "use",
		Aliases: []string{"u"},
		Short:   "Downloads and switches to the specified version.",
		Long:    "n u <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			nv := vctrl.NewNodeVersion()
			nv.UseVersion(args[0])
		},
	}
	nodeCmd.AddCommand(useCmd)

	localCmd := &cobra.Command{
		Use:     "local",
		Aliases: []string{"l"},
		Short:   "Shows installed versions.",
		Run: func(cmd *cobra.Command, args []string) {
			nv := vctrl.NewNodeVersion()
			nv.ShowInstalled()
		},
	}
	nodeCmd.AddCommand(localCmd)

	removeAllCmd := &cobra.Command{
		Use:     "remove-unused",
		Aliases: []string{"ru"},
		Short:   "Removes installed versions except the one currently in use.",
		Run: func(cmd *cobra.Command, args []string) {
			nv := vctrl.NewNodeVersion()
			nv.RemoveVersion("all")
		},
	}
	nodeCmd.AddCommand(removeAllCmd)

	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Removes the specified version.",
		Long:    "Example: n rm <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			nv := vctrl.NewNodeVersion()
			nv.RemoveVersion(args[0])
		},
	}
	nodeCmd.AddCommand(removeCmd)
	reg.Register(nodeCmd)
}
