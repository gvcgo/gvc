package cmd

import (
	"github.com/gvcgo/gvc/pkg/browser"
	"github.com/spf13/cobra"
)

func RegisterBrowser(cli *Cli) {
	parent := &cobra.Command{
		Use:     "browser",
		Aliases: []string{"b"},
		GroupID: cli.groupID,
		Short:   "Handles data from browser.",
	}

	listBrowsers := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Shows supported browsers.",
		Run: func(cmd *cobra.Command, args []string) {
			browser.ShowSupportedBrowser()
		},
	}
	parent.AddCommand(listBrowsers)

	save := &cobra.Command{
		Use:     "save",
		Aliases: []string{"s"},
		Short:   "Saves data from a local browser.",
		Long:    "g b s <browser-name>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			keep, _ := cmd.Flags().GetBool("keep-temp-files")
			browser.SaveBrowserData(args[0], keep)
		},
	}
	save.Flags().BoolP("keep-temp-files", "k", false, "Keeps temp files or not")
	parent.AddCommand(save)

	cli.rootCmd.AddCommand(parent)
}
