package langs

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

/*
Java
Maven
Gradle
*/
func SetJava(reg IRegister) {
	javaCmd := &cobra.Command{
		Use:     "java",
		Aliases: []string{"j"},
		Short:   "Java related CLIs.",
	}

	remoteCmd := &cobra.Command{
		Use:     "remote",
		Aliases: []string{"r"},
		Short:   "Shows available versions of jdk from remote websites.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewJDKVersion()
			gv.ShowVersions()
		},
	}
	javaCmd.AddCommand(remoteCmd)

	useCmd := &cobra.Command{
		Use:     "use",
		Aliases: []string{"u"},
		Short:   "Downloads and switches to the specified version of jdk.",
		Long:    "Example: j u jdk17-lts",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewJDKVersion()
			gv.UseVersion(args[0])
		},
	}
	javaCmd.AddCommand(useCmd)

	localCmd := &cobra.Command{
		Use:     "local",
		Aliases: []string{"l"},
		Short:   "Shows installed versions of jdk.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewJDKVersion()
			gv.ShowInstalled()
		},
	}
	javaCmd.AddCommand(localCmd)

	removeAllCmd := &cobra.Command{
		Use:     "remove-unused",
		Aliases: []string{"ru"},
		Short:   "Removes installed versions of jdk except the one currently in use.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewJDKVersion()
			gv.RemoveUnused()
		},
	}
	javaCmd.AddCommand(removeAllCmd)

	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "R"},
		Short:   "Removes a specified version of jdk.",
		Long:    "Example: j R jdk17-lts",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewJDKVersion()
			gv.RemoveVersion(args[0])
		},
	}
	javaCmd.AddCommand(removeCmd)

	/*
		maven related
	*/
	mRemoteCmd := &cobra.Command{
		Use:     "maven-remote",
		Aliases: []string{"mr"},
		Short:   "Shows available versions of maven from remote websites.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewMavenVersion()
			gv.ShowVersions()
		},
	}
	javaCmd.AddCommand(mRemoteCmd)

	mUseCmd := &cobra.Command{
		Use:     "maven-use",
		Aliases: []string{"mu"},
		Short:   "Downloads and switches to the specified version fo maven.",
		Long:    "Example: j mu <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewMavenVersion()
			gv.UseVersion(args[0])
		},
	}
	javaCmd.AddCommand(mUseCmd)

	mLocalCmd := &cobra.Command{
		Use:     "maven-local",
		Aliases: []string{"ml"},
		Short:   "Show installed versions of maven.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewMavenVersion()
			gv.ShowInstalled()
		},
	}
	javaCmd.AddCommand(mLocalCmd)

	mMirrorCmd := &cobra.Command{
		Use:     "maven-mirror",
		Aliases: []string{"mm"},
		Short:   "Set mirrors and local repository path for maven.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewMavenVersion()
			gv.GenSettingsFile()
		},
	}
	javaCmd.AddCommand(mMirrorCmd)

	mRemoveAllCmd := &cobra.Command{
		Use:     "maven-remove-unused",
		Aliases: []string{"mru"},
		Short:   "Removes installed versions of maven except the one currently in use.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewMavenVersion()
			gv.RemoveUnused()
		},
	}
	javaCmd.AddCommand(mRemoveAllCmd)

	mRemoveCmd := &cobra.Command{
		Use:     "maven-remove",
		Aliases: []string{"mrm", "mR"},
		Short:   "Removes a specified version fo maven.",
		Long:    "Example: j mR <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewMavenVersion()
			gv.RemoveVersion(args[0])
		},
	}
	javaCmd.AddCommand(mRemoveCmd)

	/*
		gradle related
	*/
	gRemoteCmd := &cobra.Command{
		Use:     "gradle-remote",
		Aliases: []string{"gr"},
		Short:   "Shows available versions of gradle from remote website.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGradleVersion()
			gv.ShowVersions()
		},
	}
	javaCmd.AddCommand(gRemoteCmd)

	gUseCmd := &cobra.Command{
		Use:     "gradle-use",
		Aliases: []string{"gu"},
		Short:   "Downloads and switches to the specified version fo gradle.",
		Long:    "Example: j gu <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewGradleVersion()
			gv.UseVersion(args[0])
		},
	}
	javaCmd.AddCommand(gUseCmd)

	gLocalCmd := &cobra.Command{
		Long:    "gradle-local",
		Aliases: []string{"gl"},
		Short:   "Shows installed versions fo gradle.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGradleVersion()
			gv.ShowInstalled()
		},
	}
	javaCmd.AddCommand(gLocalCmd)

	gMirrorCmd := &cobra.Command{
		Use:     "gradle-mirror",
		Aliases: []string{"gm"},
		Short:   "Set aliyun repository for gradle.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGradleVersion()
			gv.GenInitFile()
		},
	}
	javaCmd.AddCommand(gMirrorCmd)

	gRemoveAllCmd := &cobra.Command{
		Use:     "gradle-remove-unused",
		Aliases: []string{"gru"},
		Short:   "Removes installed versions of gradle except the one currently in use.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGradleVersion()
			gv.RemoveUnused()
		},
	}
	javaCmd.AddCommand(gRemoveAllCmd)

	gRemoveCmd := &cobra.Command{
		Use:     "gradle-remove",
		Aliases: []string{"grm", "gR"},
		Short:   "Removes a specified version fo gradle.",
		Long:    "Example: j gR <version>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewGradleVersion()
			gv.RemoveVersion(args[0])
		},
	}
	javaCmd.AddCommand(gRemoveCmd)

	reg.Register(javaCmd)
}
