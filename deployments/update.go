package deployments

import (
	"fmt"
	"github.com/DonDebonair/ploy/engine"
	"github.com/DonDebonair/ploy/utils"
	"github.com/spf13/cobra"
)

func Update(_ *cobra.Command, args []string) {
	deploymentsConfigPath := args[0]
	nrArgs := len(args)
	serviceIds := args[1 : nrArgs-1]
	version := args[nrArgs-1]
	deployments, err := LoadDeploymentsFromFile(deploymentsConfigPath)
	cobra.CheckErr(err)
	for _, serviceId := range serviceIds {
		service := utils.Find(deployments, func(d engine.Deployment) bool {
			return d.Id() == serviceId
		})
		if service == nil {
			cobra.CheckErr(fmt.Errorf("there is no service with id '%s'", serviceId))
		}
		(*service).SetVersion(version)
	}
	err = WriteDeploymentsToFile(deploymentsConfigPath, deployments)
	cobra.CheckErr(err)
}
