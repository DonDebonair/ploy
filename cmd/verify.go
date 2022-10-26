package cmd

import (
	"github.com/DandyDev/ploy/deployments"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify [config-file]",
	Short: "Verify deployed versions of services/applications in the specified deployment config",
	Long: `Verify deployed versions of services/applications in the specified deployment config

Given the specified deployment configuration file in yaml format, 
check which versions of the services/applications in the config file are 
currently deployed and verify if these versions match the versions 
specified in the config file.`,
	Run:  deployments.Verify,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().BoolVarP(&deployments.FailOnVersionMismatch, "fail", "f", false, "Fail if any deployments don't match expected version")
}
