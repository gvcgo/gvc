package langs

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func SetZig(reg IRegister) {
	zigCmd := &cobra.Command{
		Use:     "zig",
		Aliases: []string{"z"},
		Short:   "Zig related CLIs.",
	}

	var overwriteFlagName = "force"
	installCmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs the latest stable version of zig.",
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool(overwriteFlagName)
			v := vctrl.NewZig()
			v.Install(force)
		},
	}
	installCmd.Flags().BoolP(overwriteFlagName, "f", false, "To overwrite the old version.")
	zigCmd.AddCommand(installCmd)

	envCmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"e"},
		Short:   "Sets envs for zig.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewZig()
			v.CheckAndInitEnv()
		},
	}
	zigCmd.AddCommand(envCmd)

	analyzerCmd := &cobra.Command{
		Use:     "install-zls",
		Aliases: []string{"iz"},
		Short:   "Installs zls.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: zls
		},
	}
	zigCmd.AddCommand(analyzerCmd)

	reg.Register(zigCmd)
}
