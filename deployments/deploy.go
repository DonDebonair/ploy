package deployments

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Deploy(_ *cobra.Command, args []string) {
	deploymentsConfigPath := args[0]
	deploymentsConfig, err := LoadDeploymentsFromFile(deploymentsConfigPath)
	cobra.CheckErr(err)
	fmt.Printf("Config: %s", deploymentsConfig)
	//for _, deployment := range deploymentsConfig.Deployments {
	//	fmt.Printf("Deploying %s\n", deployment.Id)
	//}
}
