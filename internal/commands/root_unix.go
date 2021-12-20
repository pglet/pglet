//go:build !windows

package commands

import (
	"github.com/spf13/cobra"
)

var (
	version  = "unknown"
	commit   = "unknown"
	LogLevel string
)

func NewRootCmd() *cobra.Command {
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
		newServerCommand(),
		newClientCommand(),
	)

	return cmd
}
