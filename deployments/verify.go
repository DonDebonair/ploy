package deployments

import (
	"github.com/DandyDev/ploy/engine"
	"github.com/spf13/cobra"
	"sync"
)

// TODO: should this command return a non-zero exit code when the versions of one or more deployments don't match?
func Verify(_ *cobra.Command, args []string) {
	deployments, err := LoadDeploymentsFromFile(args)
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
	p := CreateDeploymentPrinter(deploymentConfig.Id())
	deploymentEngine := engine.GetEngine(deploymentConfig.Type())
	version, err := deploymentEngine.CheckVersion(deploymentConfig)
	if err != nil {
		return err
	}
	if version != deploymentConfig.Version() {
		p("❌ version '%s' does not match expected version '%s'", version, deploymentConfig.Version())
	} else {
		p("✅ version '%s' matches expected version '%s'", version, deploymentConfig.Version())
	}
	return nil
}
