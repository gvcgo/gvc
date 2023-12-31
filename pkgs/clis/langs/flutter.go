package langs

import "github.com/spf13/cobra"

func SetFlutter(reg IRegister) {
	nodeCmd := &cobra.Command{
		Use:     "flutter",
		Aliases: []string{"f"},
		Short:   "Flutter related CLIs.",
	}

	remoteCmd := &cobra.Command{}
	nodeCmd.AddCommand(remoteCmd)

	useCmd := &cobra.Command{}
	nodeCmd.AddCommand(useCmd)

	localCmd := &cobra.Command{}
	nodeCmd.AddCommand(localCmd)

	removeAllCmd := &cobra.Command{}
	nodeCmd.AddCommand(removeAllCmd)

	removeCmd := &cobra.Command{}
	nodeCmd.AddCommand(removeCmd)
	reg.Register(nodeCmd)
}
