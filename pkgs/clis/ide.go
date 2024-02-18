package clis

import (
	"github.com/gvcgo/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func (that *Cli) vscode() {
	codeCmd := &cobra.Command{
		Use:     "vscode",
		Aliases: []string{"vs"},
		Short:   "Installs vscode, extensions, etc.",
		GroupID: that.groupID,
	}

	codeCmd.AddCommand(&cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs vscode.",
		Run: func(cmd *cobra.Command, args []string) {
			gcode := vctrl.NewCode()
			gcode.Install()
		},
	})

	codeCmd.AddCommand(&cobra.Command{
		Use:     "upload-configs",
		Aliases: []string{"u"},
		Short:   "Uploads settings.json, keybindings.json, extensions.txt to remote repo.",
		Run: func(cmd *cobra.Command, args []string) {
			gcode := vctrl.NewCode()
			gcode.HandleVSCodeFiles(false)
		},
	})

	codeCmd.AddCommand(&cobra.Command{
		Use:     "download-configs",
		Aliases: []string{"d"},
		Short:   "Downloads settings.json, keybindings.json, extensions.txt from remote repo.",
		Run: func(cmd *cobra.Command, args []string) {
			gcode := vctrl.NewCode()
			gcode.HandleVSCodeFiles(true)
		},
	})

	codeCmd.AddCommand(&cobra.Command{
		Use:     "fixgit",
		Aliases: []string{"f"},
		Short:   "Fixes git.exe in Msys2/Cygwin for vscode.",
		Run: func(cmd *cobra.Command, args []string) {
			gcode := vctrl.NewCppManager()
			gcode.RepairGitForVSCode()
		},
	})
	that.rootCmd.AddCommand(codeCmd)
}

// TODO: new command.
func (that *Cli) neovim() {
	nvimCmd := &cobra.Command{
		Use:     "nvim",
		Aliases: []string{"nv"},
		Short:   "Installs neovim.",
		GroupID: that.groupID,
	}

	nvimCmd.AddCommand(&cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs neovim.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewNVim()
			v.Install()
		},
	})

	that.rootCmd.AddCommand(nvimCmd)
}
