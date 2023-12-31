package langs

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func SetVlang(reg IRegister) {
	vlangCmd := &cobra.Command{}

	var overwriteFlagName = "force"
	installCmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs the latest version of vlang.",
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool(overwriteFlagName)
			v := vctrl.NewVlang()
			v.Install(force)
		},
	}
	installCmd.Flags().BoolP(overwriteFlagName, "f", false, "To overwrite the old version.")
	vlangCmd.AddCommand(installCmd)

	envCmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"e"},
		Short:   "Sets envs for vlang.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewVlang()
			v.CheckAndInitEnv()
		},
	}
	vlangCmd.AddCommand(envCmd)

	analyzerCmd := &cobra.Command{
		Use:     "install-analyzer",
		Aliases: []string{"ia"},
		Short:   "Installs v-analyzer and its extension for VSCode.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewVlang()
			v.InstallVAnalyzerForVscode()
		},
	}
	vlangCmd.AddCommand(analyzerCmd)

	reg.Register(vlangCmd)
}
