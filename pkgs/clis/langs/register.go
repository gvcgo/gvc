package langs

import (
	"github.com/spf13/cobra"
)

type IRegister interface {
	Register(cmd *cobra.Command)
}
