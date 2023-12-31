package langs

import "github.com/spf13/cobra"

func SetJulia(reg IRegister) {
	nodeCmd := &cobra.Command{
		Use:     "julia",
		Aliases: []string{"jl", "J"},
		Short:   "Julia related CLIs.",
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
