package deployments

import (
	"fmt"
	"github.com/DonDebonair/ploy/engine"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"sync"
)

var FailOnVersionMismatch bool

func Verify(_ *cobra.Command, args []string) {
	deployments, err := LoadDeploymentsFromFile(args[0])
	cobra.CheckErr(err)

	errorChan := make(chan error, len(deployments))
	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		go func(deployment engine.Deployment) {
			defer wg.Done()
			errorChan <- verifyDeployment(deployment, FailOnVersionMismatch)
		}(deployment)
	}
	wg.Wait()
	close(errorChan)
	var result *multierror.Error
	for err := range errorChan {
		result = multierror.Append(result, err)
	}
	cobra.CheckErr(result.ErrorOrNil())
}

func verifyDeployment(deploymentConfig engine.Deployment, failOnVersionMismatch bool) error {
	p := CreateDeploymentPrinter(deploymentConfig.Id())
	deploymentEngine, err := engine.GetEngine(deploymentConfig.Type())
	if err != nil {
		return err
	}
	version, err := deploymentEngine.CheckVersion(deploymentConfig)
	if err != nil {
		return err
	}
	if version != deploymentConfig.Version() {
		p("❌ version '%s' does not match expected version '%s'", version, deploymentConfig.Version())
		if failOnVersionMismatch {
			return fmt.Errorf("%s: version '%s' does not match expected version '%s'", deploymentConfig.Id(), version, deploymentConfig.Version())
		}
	} else {
		p("✅ version '%s' matches expected version '%s'", version, deploymentConfig.Version())
	}
	return nil
}
