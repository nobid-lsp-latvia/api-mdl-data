// SPDX-License-Identifier: EUPL-1.2

package main

import (
	"github.com/spf13/cobra"
)

// mdlCmd represents the health command.
var mdlCmd = &cobra.Command{
	Use:           "health",
	Short:         "Check health of the server",
	Long:          `Check if the web server is running and responding to healthz request`,
	RunE:          runMdl,
	SilenceErrors: true,
}

func runMdl(cmd *cobra.Command, _ []string) error {
	checker := &HealthCheck{
		Path:    "/mdl",
		Version: Version,
	}

	return checker.RunHealthCheck(cmd)
}

func init() {
	initRootCmd()
	RootCmd.AddCommand(mdlCmd)
}
