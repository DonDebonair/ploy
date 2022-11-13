package deployments

import (
	"fmt"
	"github.com/DonDebonair/ploy/engine"
	"github.com/spf13/cobra"
)

func Update(_ *cobra.Command, args []string) {
	deploymentsConfigPath := args[0]
	serviceId := args[1]
	version := args[2]
	deployments, err := LoadDeploymentsFromFile(deploymentsConfigPath)
	cobra.CheckErr(err)
	service := Find(deployments, func(d engine.Deployment) bool {
		return d.Id() == serviceId
	})
	if service == nil {
		cobra.CheckErr(fmt.Errorf("there is no service with id '%s'", serviceId))
	}
	(*service).SetVersion(version)
	err = WriteDeploymentsToFile(deploymentsConfigPath, deployments)
	cobra.CheckErr(err)
}
