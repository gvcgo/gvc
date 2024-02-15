package langs

import (
	"github.com/gvcgo/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func SetFlutter(reg IRegister) {
	flutterCmd := &cobra.Command{
		Use:     "flutter",
		Aliases: []string{"f"},
		Short:   "Flutter related CLIs.",
	}

	remoteCmd := &cobra.Command{
		Use:     "remote",
		Aliases: []string{"r"},
		Short:   "Shows available versions from remote website.",
		Run: func(cmd *cobra.Command, args []string) {
			fv := vctrl.NewFlutterVersion()
			fv.ShowVersions()
		},
	}
	flutterCmd.AddCommand(remoteCmd)

	useCmd := &cobra.Command{
		Use:     "use",
		Aliases: []string{"u"},
		Short:   "Downloads and switches to the specified version.",
		Long:    "f u <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			fv := vctrl.NewFlutterVersion()
			fv.UseVersion(args[0])
		},
	}
	flutterCmd.AddCommand(useCmd)

	envCmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"e"},
		Short:   "Sets envs for flutter.",
		Run: func(cmd *cobra.Command, args []string) {
			fv := vctrl.NewFlutterVersion()
			fv.CheckAndInitEnv()
		},
	}
	flutterCmd.AddCommand(envCmd)

	localCmd := &cobra.Command{
		Use:     "local",
		Aliases: []string{"l"},
		Short:   "Shows installed versions.",
		Run: func(cmd *cobra.Command, args []string) {
			fv := vctrl.NewFlutterVersion()
			fv.ShowInstalled()
		},
	}
	flutterCmd.AddCommand(localCmd)

	removeAllCmd := &cobra.Command{
		Use:     "remove-unused",
		Aliases: []string{"ru"},
		Short:   "Removes installed versions except the one currently in use.",
		Run: func(cmd *cobra.Command, args []string) {
			fv := vctrl.NewFlutterVersion()
			fv.RemoveUnused()
		},
	}
	flutterCmd.AddCommand(removeAllCmd)

	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Removes a specified version.",
		Long:    "Examples: f rm <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			fv := vctrl.NewFlutterVersion()
			fv.RemoveVersion(args[0])
		},
	}
	flutterCmd.AddCommand(removeCmd)

	androidCmd := &cobra.Command{
		Use:     "android-sdkmanager",
		Aliases: []string{"as"},
		Short:   "Installs android cmdline tools(sdkmanager).",
		Run: func(cmd *cobra.Command, args []string) {
			fv := vctrl.NewFlutterVersion()
			fv.InstallAndroidTool()
		},
	}
	flutterCmd.AddCommand(androidCmd)

	avdCmd := &cobra.Command{
		Use:     "tools-avd",
		Aliases: []string{"ta"},
		Short:   "Install build-tools, platform-tools, etc. And create avd for android.",
		Long:    "Example: f ta <your-avd-name>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			fv := vctrl.NewFlutterVersion()
			fv.SetupAVD(args[0])
		},
	}
	flutterCmd.AddCommand(avdCmd)

	startAvdCmd := &cobra.Command{
		Use:     "start-avd",
		Aliases: []string{"savd", "sa"},
		Short:   "Start an avd.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewFlutterVersion()
			gv.StartAVD()
		},
	}
	flutterCmd.AddCommand(startAvdCmd)

	fixRepoForFlutterProjectCmd := &cobra.Command{
		Use:     "aliyun-repo",
		Aliases: []string{"arepo", "ar"},
		Short:   "Use aliyun repo for a flutter project.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewFlutterVersion()
			gv.ReplaceMavenRepo()
		},
	}
	flutterCmd.AddCommand(fixRepoForFlutterProjectCmd)

	reg.Register(flutterCmd)
}
