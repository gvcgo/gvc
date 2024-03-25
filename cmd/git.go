package cmd

import (
	"github.com/gvcgo/gvc/pkg/git"
	"github.com/spf13/cobra"
)

func RegisterGit(cli *Cli) {
	parent := &cobra.Command{
		Use:     "git",
		Aliases: []string{"g"},
		Short:   "Git related CLIs.",
		GroupID: cli.groupID,
	}

	hosts := &cobra.Command{
		Use:     "update-hosts",
		Aliases: []string{"uh", "u"},
		Short:   "Updates hosts file.",
		Run: func(cmd *cobra.Command, args []string) {
			m := git.NewModifier()
			m.Run()
		},
	}
	parent.AddCommand(hosts)

	var (
		destHostName string = "dest_host"
		destPortName string = "dest_port"
		timeoutName  string = "timeout"
	)
	crokscrew := &cobra.Command{
		Use:     "crokscrew",
		Aliases: []string{"cs", "c"},
		Short:   "Http proxy for ssh.",
		Long:    "Example: g g cs --dest_host=xxx --dest_port=xxx --timeout=xxx",
		Run: func(cmd *cobra.Command, args []string) {
			destHost, _ := cmd.Flags().GetString(destHostName)
			destPort, _ := cmd.Flags().GetString(destPortName)
			timeout, _ := cmd.Flags().GetInt(timeoutName)
			if destHost == "" || destPort == "" {
				cmd.Help()
				return
			}
			git.GrokscrewHttpSSH(destHost, destPort, timeout)
		},
	}
	crokscrew.Flags().StringP(destHostName, "a", "", "Specifies dest host.")
	crokscrew.Flags().StringP(destPortName, "p", "", "Specifies dest port.")
	crokscrew.Flags().IntP(timeoutName, "t", 10, "Specifies timeout.")
	parent.AddCommand(crokscrew)

	toggle := &cobra.Command{
		Use:     "toggle-proxy",
		Aliases: []string{"tp", "t"},
		Short:   "Toggle proxy for git ssh.",
		Run: func(cmd *cobra.Command, args []string) {
			git.ToggleProxyForSSH()
		},
	}
	parent.AddCommand(toggle)

	cli.rootCmd.AddCommand(parent)
}
