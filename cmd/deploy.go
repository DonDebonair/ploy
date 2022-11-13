package cmd

import (
	"github.com/DonDebonair/ploy/deployments"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy <config-file>",
	Short: "Deploy services/applications in the specified deployment config",
	Long: `Deploy services/applications in the specified deployment config

Given the specified deployment configuration file in yaml format, 
check which versions of the services/applications in the config file are 
currently deployed and if these differ from the versions specified in the 
configuration file, deploy the specified versions.`,
	Run:     deployments.Deploy,
	Args:    cobra.ExactArgs(1),
	Example: "ploy deploy production.yml",
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
