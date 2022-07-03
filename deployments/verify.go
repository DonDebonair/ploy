package deployments

import (
	"fmt"
	"github.com/DandyDev/ploy/engine"
	"github.com/spf13/cobra"
	"sync"
)

// TODO: should this command return a non-zero exit code when the versions of one or more deployments don't match?
func Verify(_ *cobra.Command, args []string) {
	deploymentsConfigPath := args[0]
	deployments, err := LoadDeploymentsFromFile(deploymentsConfigPath)
	cobra.CheckErr(err)
	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		deployment := deployment
		go func() {
			defer wg.Done()
			err := verifyDeployment(deployment)
			cobra.CheckErr(err)
		}()
	}
	wg.Wait()
}

func verifyDeployment(deploymentConfig engine.Deployment) error {
	deploymentEngine := engine.GetEngine(deploymentConfig.Type())
	version, err := deploymentEngine.CheckVersion(deploymentConfig)
	if err != nil {
		return err
	}
	if version != deploymentConfig.Version() {
		fmt.Printf("❌ Deployment %s version '%s' does not match expected version '%s'\n", deploymentConfig.Id(), version, deploymentConfig.Version())
	} else {
		fmt.Printf("✅ Deployment %s version '%s' matches expected version '%s'\n", deploymentConfig.Id(), version, deploymentConfig.Version())
	}
	return nil
}
