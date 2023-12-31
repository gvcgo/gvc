package langs

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func SetTypst(reg IRegister) {
	typstCmd := &cobra.Command{}

	var overWriteFlagName = "force"
	installCmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs the latest version of typst.",
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool(overWriteFlagName)
			v := vctrl.NewTypstVersion()
			v.Install(force)
		},
	}
	installCmd.Flags().BoolP(overWriteFlagName, "f", false, "To overwirte the old version.")
	typstCmd.AddCommand(installCmd)

	envCmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"e"},
		Short:   "Set envs for typst.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewTypstVersion()
			v.CheckAndInitEnv()
		},
	}
	typstCmd.AddCommand(envCmd)

	reg.Register(typstCmd)
}
