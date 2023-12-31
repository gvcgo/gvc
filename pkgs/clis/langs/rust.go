package langs

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func SetRust(reg IRegister) {
	rustCmd := &cobra.Command{
		Use:     "rust",
		Aliases: []string{"r"},
	}

	installCmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Install the latest rust compiler tools.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewRustInstaller()
			v.Install()
		},
	}
	rustCmd.AddCommand(installCmd)

	envCmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"e"},
		Short:   "Set acceleration envs for rust in China.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewRustInstaller()
			v.SetAccelerationEnv()
		},
	}
	rustCmd.AddCommand(envCmd)

	reg.Register(rustCmd)
}
