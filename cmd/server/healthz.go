// SPDX-License-Identifier: EUPL-1.2

package main

import (
	"github.com/spf13/cobra"
)

// healthCmd represents the health command.
var healthCmd = &cobra.Command{
	Use:           "health",
	Short:         "Check health of the server",
	Long:          `Check if the web server is running and responding to healthz request`,
	RunE:          runHealth,
	SilenceErrors: true,
}

func runHealth(cmd *cobra.Command, _ []string) error {
	checker := &HealthCheck{
		Path:    "/healthz",
		Version: Version,
	}

	return checker.RunHealthCheck(cmd)
}

func init() {
	initRootCmd()
	RootCmd.AddCommand(healthCmd)
}
