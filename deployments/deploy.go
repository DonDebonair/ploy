package deployments

import (
	"github.com/DandyDev/ploy/engine"
	"github.com/spf13/cobra"
	"sync"
)

func Deploy(_ *cobra.Command, args []string) {
	deployments, err := LoadDeploymentsFromFile(args)
	cobra.CheckErr(err)
	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		deployment := deployment
		go func() {
			defer wg.Done()
			err := doDeployment(deployment)
			cobra.CheckErr(err)
		}()
	}
	wg.Wait()
}

func doDeployment(deploymentConfig engine.Deployment) error {
	p := CreateDeploymentPrinter(deploymentConfig.Id())
	deploymentEngine := engine.GetEngine(deploymentConfig.Type())
	p("checking deployed version...")
	version, err := deploymentEngine.CheckVersion(deploymentConfig)
	if err != nil {
		return err
	}
	if version != deploymentConfig.Version() {
		p("version '%s' does not match expected version '%s'. Deploying new version...", version, deploymentConfig.Version())
		err = deploymentEngine.Deploy(deploymentConfig, p)
		if err != nil {
			return err
		} else {
			p("version %s deployed successfully!", deploymentConfig.Version())
			return nil
		}
	} else {
		p("version '%s' matches expected version '%s'. Skipping...", version, deploymentConfig.Version())
		return nil
	}
}
