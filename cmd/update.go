package cmd

import (
	"github.com/DonDebonair/ploy/deployments"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
//
//goland:noinspection SqlNoDataSourceInspection
var updateCmd = &cobra.Command{
	Use:   "update <config-file> <service-id> <version>",
	Short: "Update the version of a service",
	Long: `Update the specified service to the desired version

This will update the yaml file so that the specified service is now set to
the desired version.`,
	Run:  deployments.Update,
	Args: cobra.ExactArgs(3),
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
