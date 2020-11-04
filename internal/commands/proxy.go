package commands

import (
	"github.com/pglet/pglet/internal/proxy"
	"github.com/spf13/cobra"
)

func newProxyCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "proxy",
		Short: "Start proxy service",
		Long:  `Proxy is for ...`,
		Run: func(cmd *cobra.Command, args []string) {
			proxy.RunService(cmd.Context())
		},
	}

	return cmd
}
