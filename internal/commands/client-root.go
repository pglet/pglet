//go:build !windows

package commands

import (
	"github.com/spf13/cobra"
)

func NewClientRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pglet",
		Short:   "Pglet",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configLogging()
		},
	}

	cmd.SetVersionTemplate("{{.Version}}")

	cmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "verbosity level for logs")

	cmd.AddCommand(
		newPageCommand(),
		newAppCommand(),
		newClientCommand(),
	)

	return cmd
}
