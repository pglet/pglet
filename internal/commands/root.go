package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "unknown"
	commit  = "unknown"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pglet",
		Short:   "Pglet",
		Version: fmt.Sprintf("%s-%s", version, commit),
	}

	cmd.SetVersionTemplate("{{.Version}}")

	cmd.AddCommand(
		newPageCommand(),
		newAppCommand(),
		newProxyCommand(),
		newServerCommand(),
	)

	return cmd
}
