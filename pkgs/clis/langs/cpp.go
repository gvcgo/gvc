package langs

import (
	"github.com/gvcgo/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func SetCpp(reg IRegister) {
	cppCmd := cobra.Command{
		Use:     "cpp",
		Aliases: []string{"cp"},
		Short:   "Cpp related CLIs.",
	}

	msys2InstallCmd := &cobra.Command{
		Use:     "install-msys2",
		Aliases: []string{"im"},
		Short:   "Installs the latest version of msys2.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewCppManager()
			v.InstallMsys2()
		},
	}
	cppCmd.AddCommand(msys2InstallCmd)

	uninstallMsys2Cmd := &cobra.Command{
		Use:     "uninstall-msys2",
		Aliases: []string{"um"},
		Short:   "Uninstalls msys2.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewCppManager()
			v.UninstallMsys2()
		},
	}
	cppCmd.AddCommand(uninstallMsys2Cmd)

	cygwinInstallCmd := &cobra.Command{
		Use:     "install-cygwin",
		Aliases: []string{"ic"},
		Short:   "Installs cygwin.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewCppManager()
			v.InstallCygwin("")
		},
	}
	cppCmd.AddCommand(cygwinInstallCmd)

	vcpkgInstallCmd := &cobra.Command{
		Use:     "install-vcpkg",
		Aliases: []string{"iv"},
		Short:   "Installs vcpkg.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewCppManager()
			v.InstallVCPkg()
		},
	}
	cppCmd.AddCommand(vcpkgInstallCmd)

	reg.Register(&cppCmd)
}
