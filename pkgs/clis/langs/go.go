package langs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/gvc/pkgs/utils"
	"github.com/gvcgo/gvc/pkgs/utils/sorts"
	"github.com/gvcgo/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

const (
	TempChar = "@"
)

func RecoverArgs(args ...string) (aList []string) {
	for _, v := range args {
		if strings.HasPrefix(v, TempChar) {
			aList = append(aList, strings.Replace(v, TempChar, "-", 1))
		} else {
			aList = append(aList, v)
		}
	}
	return aList
}

/*
Commad list:
1. show remote versions.
2. download and use a specified version.
3. list installed versions.
4. remove installed versions except the one currently in use.
5. remove a specified version.
6. set envs for go.
7. list the distributions for go.
8. search packages written in go.
9. build go project for diffrent platforms. Just script free.
10. automatically rename a local package.
11. install protoc.
12. install proto-gen-go.
13. install protoc-gen-go-grpc.
*/
func SetGo(reg IRegister) {
	goCmd := &cobra.Command{
		Use:     "go",
		Aliases: []string{"G"},
		Short:   "Golang related CLIs.",
	}

	var showAllFlagName = "all"
	remoteCmd := &cobra.Command{
		Use:     "remote",
		Aliases: []string{"r"},
		Short:   "Shows available versions from remote website.",
		Run: func(cmd *cobra.Command, args []string) {
			showall, _ := cmd.Flags().GetBool(showAllFlagName)
			gv := vctrl.NewGoVersion()
			arg := vctrl.ShowStable
			if showall {
				arg = vctrl.ShowAll
			}
			gv.ShowRemoteVersions(arg)
		},
	}
	remoteCmd.Flags().BoolP(showAllFlagName, "a", false, "To show all remote versions.")
	goCmd.AddCommand(remoteCmd)

	useCmd := &cobra.Command{
		Use:     "use",
		Aliases: []string{"u"},
		Short:   "Downloads and switches to the specified version.",
		Long:    "Example: G u 1.21.5",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewGoVersion()
			gv.UseVersion(args[0])
		},
	}
	goCmd.AddCommand(useCmd)

	localCmd := &cobra.Command{
		Use:     "local",
		Aliases: []string{"l"},
		Short:   "Shows all installed versions.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGoVersion()
			gv.ShowInstalled()
		},
	}
	goCmd.AddCommand(localCmd)

	removeAllCmd := &cobra.Command{
		Use:     "remove-unused",
		Aliases: []string{"R"},
		Short:   "Remove installed versions except the one currently in use.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGoVersion()
			gv.RemoveUnused()
		},
	}
	goCmd.AddCommand(removeAllCmd)

	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove an installed version.",
		Long:    "Example: G rm 1.21.5",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewGoVersion()
			gv.RemoveVersion(args[0])
		},
	}
	goCmd.AddCommand(removeCmd)

	envsCmd := &cobra.Command{
		Use:     "envs",
		Aliases: []string{"e"},
		Short:   "Set envs for go, including GOPATH, GOBIN, GOPROXY, etc.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGoVersion()
			gv.CheckAndInitEnv()
		},
	}
	goCmd.AddCommand(envsCmd)

	distCmd := &cobra.Command{
		Use:     "list-distros",
		Aliases: []string{"ld"},
		Short:   "List the platforms supported by go compilers.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGoVersion()
			gv.ShowGoDistlist()
		},
	}
	goCmd.AddCommand(distCmd)

	var (
		orderByTimeFlag string = "order-by-time"
	)
	searchPkgCmd := &cobra.Command{
		Use:     "search",
		Aliases: []string{"s"},
		Short:   "Search for go packages.",
		Long:    "Example: G s <package-name-keyword>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			var orderBy int = sorts.ByImported
			if byTime, _ := cmd.Flags().GetBool(orderByTimeFlag); byTime {
				orderBy = sorts.ByUpdate
			}
			gv := vctrl.NewGoVersion()
			gv.SearchLibs(args[0], orderBy)
		},
	}
	searchPkgCmd.Flags().BoolP(orderByTimeFlag, "t", false, "Sort result by update update time.")
	goCmd.AddCommand(searchPkgCmd)

	buildCmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   `Compiles go code for multi-platforms [with <-ldflags "-s -w"> builtin].`,
		Long:    `If you are planning to use "-X", then remember to replace any "$" by "#".`,
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewGoVersion()
			if len(os.Args) > 3 {
				args = os.Args[3:]
			} else {
				args = []string{}
			}
			gv.Build(RecoverArgs(args...)...)
		},
	}
	goCmd.AddCommand(buildCmd)

	renameCmd := &cobra.Command{
		Use:     "remame-to",
		Aliases: []string{"rt"},
		Short:   "Renames a local package to a new name.",
		Long:    "Example: G rt <your-new-name>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			gv := vctrl.NewGoVersion()
			moduleDir, _ := os.Getwd()
			if ok, _ := utils.PathIsExist(filepath.Join(moduleDir, "go.mod")); !ok {
				gprint.PrintError("Can not find go.mod in current working dir.")
				return
			}
			gv.RenameLocalModule(moduleDir, args[0])
		},
	}
	goCmd.AddCommand(renameCmd)

	/*
		protobuf and grpc
	*/
	var forceToReplaceFlag string = "force"
	protoCmd := &cobra.Command{
		Use:     "protoc-install",
		Aliases: []string{"pi"},
		Short:   "Installs latest version of protoc.",
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool(forceToReplaceFlag)
			v := vctrl.NewProtobuffer()
			v.Install(force)
		},
	}
	protoCmd.Flags().BoolP(forceToReplaceFlag, "f", false, "Force to overwrite the old one.")
	goCmd.AddCommand(protoCmd)

	protoGoCmd := &cobra.Command{
		Use:     "proto-go-plugin",
		Aliases: []string{"pg"},
		Short:   "Installs proto-gen-go.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewProtobuffer()
			v.InstallGoProtobufPlugin()
		},
	}
	goCmd.AddCommand(protoGoCmd)

	grpcCmd := &cobra.Command{
		Use:     "grpc-install",
		Aliases: []string{"gi"},
		Short:   "Install protoc-gen-go-grpc.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewProtobuffer()
			v.InstallGoProtoGRPCPlugin()
		},
	}
	goCmd.AddCommand(grpcCmd)

	reg.Register(goCmd)
}
